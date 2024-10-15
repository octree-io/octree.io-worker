package workers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	openai "github.com/sashabaranov/go-openai"
)

func processTriviaGrading() {
	prompt := `I want you to grade these answers for these questions. For each question, put either a Yes or No for whether or not it passes an interview or an exam. Explain in-depth what the right answer is supposed to be. Be strict about the grading to make sure that the explanations are correct. It is acceptable if there are no specific examples unless the question specifically asks for examples. Q: What is a thread? A: A thread is another instance of a program running within the same program.`

	start := time.Now()

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)

	fmt.Printf("Response took %v to complete", time.Since(start))
}

func SpawnTriviaWorker(id int, msgs <-chan ampq.Delivery) {
	for msg := range msgs {
		log.Printf("[Trivia Worker %d] Received message: %s", id, msg.Body)

		processTriviaGrading()

		if err := msg.Ack(false); err != nil {
			log.Printf("[Trivia Worker %d] Failed to ack message: %v", id, err)
		} else {
			log.Printf("[Trivia Worker %d] Message ack'd", id)
		}
	}
}
