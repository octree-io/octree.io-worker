package workers

import (
	"encoding/json"
	"log"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	"octree.io-worker/internal/facade"
	testharness "octree.io-worker/internal/test_harness"
	"octree.io-worker/utils/converters"
)

type StdoutItem struct {
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
	Stdout                     []StdoutItem `json:"stdout"`
	Stderr                     []string     `json:"stderr"`
	Truncated                  bool         `json:"truncated"`
	ExecTime                   string       `json:"execTime"`
	ProcessExecutionResultTime float64      `json:"processExecutionResultTime"`
	DidExecute                 bool         `json:"didExecute"`
	BuildResult                BuildResult  `json:"buildResult"`
}

func processCompilationRequest() {
	start := time.Now()

	args := map[string]string{
		"nums":   "int[]",
		"target": "int",
	}

	testCases := []map[string]interface{}{
		{
			"nums":   []interface{}{2, 7, 11, 15},
			"target": 9,
			"output": []interface{}{0, 1},
		},
		{
			"nums":   []interface{}{3, 2, 4},
			"target": 6,
			"output": []interface{}{1, 2},
		},
		{
			"nums":   []interface{}{3, 3},
			"target": 6,
			"output": []interface{}{0, 1},
		},
	}

	code := `class Solution:
    def solve(self, nums: List[int], target: int) -> List[int]:
        seen = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in seen:
                return [seen[complement], i]
            seen[num] = i
        return []`

	returnType := "int[]"

	wrappedCode := testharness.PythonHarness(code, args, testCases, returnType)

	output, err := facade.CompilerExplorer("python", wrappedCode)
	if err != nil {
		log.Printf("Error while executing compile: %v", err)
	}

	var jsonOutput CompilerExplorerResponse
	json.Unmarshal(([]byte)(output), &jsonOutput)

	elapsed := time.Since(start)

	log.Printf("Request took %s to execute", elapsed)

	if jsonOutput.Code != 0 {
		log.Println(output)
	}

	if len(testCases) != len(jsonOutput.Stdout) {
		log.Println("Test cases and stdout lengths differ")
	} else {
		for i, item := range jsonOutput.Stdout {
			convertedJson, _ := converters.JsonToPython(testCases[i]["output"])
			log.Printf("Expected: %s, Actual: %s, Equivalence: %t",
				convertedJson, item.Text, convertedJson == item.Text)
		}
	}
}

func SpawnCompilationWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Compilation Worker %d] Received message: %s", id, msg.Body)

		processCompilationRequest()

		if err := msg.Ack(false); err != nil {
			log.Printf("[Compilation Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Compilation Worker %d] Message ack'd", id)
		}
	}
}
