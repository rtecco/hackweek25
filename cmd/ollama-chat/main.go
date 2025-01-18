package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/ollama/ollama/api"
)

const promptTemplate = `Based on all provided information (conversation, documents, and images) provide a set of specific descriptive tags of the specialities or experiences that the text
describes. Only include tags that you have high confidence about. Be specific instead of general. Do not include any additional text before your response.
Include any geographic locality where possible.

{{.Input}}

Response format:
RESULTS:
Specific Tags: [tag1, tag2, tag3, ...]
Locality: [region1, region2, ...]

Response:`

func main() {

	client, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.New("prompt").Parse(promptTemplate))

	// messages := []api.Message{
	// 	api.Message{
	// 		Role:    "system",
	// 		Content: prompt,
	// 	},
	// }
	// 	api.Message{
	// 		Role:    "user",
	// 		Content: "Name some unusual animals",
	// 	},
	// 	api.Message{
	// 		Role:    "assistant",
	// 		Content: "Monotreme, platypus, echidna",
	// 	},
	// 	api.Message{
	// 		Role:    "user",
	// 		Content: "which of these is the most dangerous?",
	// 	},
	// }

	fmt.Print("You ---> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var buf bytes.Buffer
	tmpl.Execute(&buf, struct{ Input string }{input})

	expandedPrompt := buf.String()
	fmt.Println(expandedPrompt)

	req := &api.GenerateRequest{
		Model:  "llama3.2-vision",
		Prompt: expandedPrompt,

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Printf("%#v\n", resp)

		return nil
	}

	err = client.Generate(ctx, req, respFunc)

	if err != nil {
		log.Fatal(err)
	}

	// messages := []api.Message{}
	// messages = append(messages, api.Message{Role: "user", Content: expandedTemplate})

	// // if input == ":done" {
	// // 	break
	// // }

	// ctx := context.Background()
	// req := &api.ChatRequest{
	// 	Model:    "llama3.2-vision",
	// 	Messages: messages,
	// }

	// respFunc := func(resp api.ChatResponse) error {
	// 	fmt.Print(resp.Message.Content)
	// 	return nil
	// }

	// err = client.Chat(ctx, req, respFunc)

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func readDocument(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
