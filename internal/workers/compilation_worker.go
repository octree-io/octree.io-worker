package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"octree.io-worker/internal/clients"
	"octree.io-worker/internal/facade"
	testharness "octree.io-worker/internal/test_harness"
	"octree.io-worker/internal/utils"
)

type OutputItem struct {
	Text string `json:"text"`
}

type BuildResult struct {
	Code               int      `json:"code"`
	TimedOut           bool     `json:"timedOut"`
	Stdout             []string `json:"stdout"`
	Stderr             []string `json:"stderr"`
	Downloads          []string `json:"downloads"`
	ExecutableFilename string   `json:"executableFilename"`
	CompilationOptions []string `json:"compilationOptions"`
}

type CompilerExplorerResponse struct {
	Code                       int          `json:"code"`
	OkToCache                  bool         `json:"okToCache"`
	TimedOut                   bool         `json:"timedOut"`
	Stdout                     []OutputItem `json:"stdout"`
	Stderr                     []OutputItem `json:"stderr"`
	Truncated                  bool         `json:"truncated"`
	ExecTime                   int          `json:"execTime"`
	ProcessExecutionResultTime float64      `json:"processExecutionResultTime"`
	DidExecute                 bool         `json:"didExecute"`
	BuildResult                BuildResult  `json:"buildResult"`
}

type CompilationRequestMessage struct {
	SubmissionId string `json:"submissionId"`
	SocketId     string `json:"socketId"`
}

type CompilationResponseMessage struct {
	SubmissionId string `json:"submissionId"`
	SocketId     string `json:"socketId"`
	RoomId       string `json:"roomId"`
	Username     string `json:"username"`
	Language     string `json:"language"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	ExecTime     string `json:"execTime"`
}

func queryProblemByID(client *mongo.Client, id int) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("octree").Collection("problems")

	filter := bson.M{"id": id}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to find problem by ID: %v", err)
	}

	return result, nil
}

func sendCompilationResponseMessage(response CompilationResponseMessage) error {
	conn, err := clients.GetRabbitMQConnection()
	if err != nil {
		return fmt.Errorf("failed to get RabbitMQ connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	queueName := "compilation_responses"
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	messageBody, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response message: %w", err)
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		ampq.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish compilation response message: %w", err)
	}

	log.Printf("Compilation response message sent to queue: %s\n", queueName)
	return nil
}

func processCompilationRequest(msg ampq.Delivery) {
	var message CompilationRequestMessage

	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Printf("Failed to parse message to JSON: %v\n", err)
		return
	}

	submissionId := message.SubmissionId
	socketId := message.SocketId

	if submissionId == "" {
		log.Println("SubmissionId is missing or empty")
		return
	}

	if socketId == "" {
		log.Println("SocketId is missing or empty")
		return
	}

	ctx := context.Background()
	pgPool, err := clients.GetPostgresPool()
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v\n", err)
	}

	var (
		problemId int
		language  string
		code      string
		runType   string
		roomId    string
		username  string
	)

	err = pgPool.QueryRow(
		ctx,
		"SELECT problem_id, language, code, type, room_id, username FROM submissions WHERE submission_id=$1", submissionId,
	).Scan(&problemId, &language, &code, &runType, &roomId, &username)

	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return
	}
	log.Printf("Problem ID: %d\nLanguage: %s\nCode: %s\nRun type: %s\nRoom ID: %s\n", problemId, language, code, runType, roomId)

	client, err := clients.GetMongoClient()
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	problem, err := queryProblemByID(client, problemId)
	if err != nil {
		log.Printf("Error finding problem: %v", err)
		return
	}

	args := utils.ConvertBsonMToStringMap(problem["args"].(bson.M))
	testCases := []map[string]interface{}{}
	outputs := []map[string]interface{}{}
	returnType := problem["returnType"].(string)

	answerAnyOrder, ok := problem["answerAnyOrder"].(bool)
	if !ok {
		answerAnyOrder = false
	}

	deepSort, ok := problem["deepSort"].(bool)
	if !ok {
		deepSort = false
	}

	log.Printf("answerAnyOrder: %v\ndeepSort: %v\n", answerAnyOrder, deepSort)

	testCasesKey := "sampleTestCases"
	if runType == "submit" {
		testCasesKey = "judgeTestCases"
	}

	for _, entry := range problem[testCasesKey].(bson.A) {
		entryBytes, err := bson.Marshal(entry)
		if err != nil {
			log.Printf("Failed to marshal bson.D: %v", err)
			return
		}

		var entryMap bson.M
		err = bson.Unmarshal(entryBytes, &entryMap)
		if err != nil {
			log.Printf("Failed to unmarshal to bson.M: %v", err)
			return
		}

		input := map[string]interface{}(entryMap["input"].(bson.M))
		output := entryMap["output"]

		outputMap := map[string]interface{}{
			"output": utils.ConvertBsonToNative(output),
		}

		testCases = append(testCases, input)
		outputs = append(outputs, outputMap)
	}

	var wrappedCode string

	switch language {
	case "python":
		wrappedCode = testharness.PythonHarness(code, args, testCases, returnType)

	case "cpp":
		wrappedCode = testharness.CppHarness(code, args, testCases, returnType)

	case "csharp":
		wrappedCode = testharness.CsharpHarness(code, args, testCases, returnType)

	case "java":
		wrappedCode = testharness.JavaHarness(code, args, testCases, returnType)

	case "ruby":
		wrappedCode = testharness.RubyHarness(code, args, testCases, returnType)

	case "javascript":
		wrappedCode = testharness.JavaScriptHarness(code, args, testCases, returnType)

	case "typescript":
		wrappedCode = testharness.TypeScriptHarness(code, args, testCases, returnType)

	default:
		fmt.Println("Unsupported language")
		return
	}

	var stdout, stderr string
	var execTime int

	switch language {
	case "typescript":
		start := time.Now()
		stdout, stderr, err = facade.ExecuteTypeScript(language, wrappedCode)
		elapsed := time.Since(start).Milliseconds()

		execTime = int(elapsed)

		if err != nil {
			fmt.Println("Error while executing TypeScript")
		}

	case "javascript":
		start := time.Now()
		stdout, stderr, err = facade.ExecuteJavaScript(language, wrappedCode)
		elapsed := time.Since(start).Milliseconds()
		execTime = int(elapsed)
		if err != nil {
			fmt.Println("Error while executing JavaScript")
		}

	default:
		start := time.Now()

		output, err := facade.CompilerExplorer(language, wrappedCode)
		if err != nil {
			log.Printf("Error while executing compile: %v", err)
		}

		var jsonOutput CompilerExplorerResponse
		json.Unmarshal(([]byte)(output), &jsonOutput)

		elapsed := time.Since(start).Milliseconds()

		execTime = int(jsonOutput.ExecTime)

		for _, out := range jsonOutput.Stdout {
			stdout += out.Text + "\n"
		}

		for _, err := range jsonOutput.Stderr {
			stderr += err.Text + "\n"
		}

		log.Printf("Request took %s to execute and %s to run", strconv.Itoa(execTime), strconv.Itoa(int(elapsed)))
	}

	fmt.Printf("Exec time: %s\n", strconv.Itoa(execTime))

	outputString := fmt.Sprintf(`{"stdout": "%s", "stderr": "%s", "execTime": %s}`, stdout, stderr, strconv.Itoa(execTime))

	status := "SUCCEEDED"
	if runType == "submit" {
		result := facade.JudgeTestCases(outputs, stdout, answerAnyOrder, deepSort, returnType)
		fmt.Printf("Verdict: %v\n", result)
		if !result {
			status = "FAILED"
		}
	}

	updateQuery := `
    UPDATE submissions
		SET output = $1, status = $2
		WHERE submission_id = $3;
  `

	_, err = pgPool.Exec(ctx, updateQuery, outputString, status, submissionId)
	if err != nil {
		log.Printf("Failed to update submission: %v\n", err)
	}

	responseMessage := CompilationResponseMessage{
		SubmissionId: submissionId,
		SocketId:     socketId,
		Username:     username,
		RoomId:       roomId,
		Language:     language,
		Type:         runType,
		Status:       status,
		Stdout:       stdout,
		Stderr:       stderr,
		ExecTime:     strconv.Itoa(execTime),
	}

	err = sendCompilationResponseMessage(responseMessage)
	if err != nil {
		log.Printf("Failed to send a compilation response message: %v", err)
	}
}

func SpawnCompilationWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Compilation Worker %d] Received message: %s", id, msg.Body)

		processCompilationRequest(msg)

		if err := msg.Ack(false); err != nil {
			log.Printf("[Compilation Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Compilation Worker %d] Message ack'd", id)
		}
	}
}
