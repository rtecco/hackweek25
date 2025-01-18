package chat

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/ollama/ollama/api"
)

var buyerPrompt = `Based on the provided information provide a search query for a professional that can help resolve the user's issue. The query should be no
more than 4 words. Also ask ONE follow up question asking for more specifics. Follow the Response Format exactly.

Current information:
{{.Input}}

Response Format:
RESULTS:
Search Query: <query>
Follow Up Question: <question>

Response:
`

var buyerPromptTemplate = template.Must(template.New("buyer").Parse(buyerPrompt))

var sellerPrompt = `Based on all provided information (conversation, documents, and images) provide a set of specific descriptive tags of the specialities or experiences that the text
		describes. Only include tags that you have high confidence about. Be specific instead of general. Do not include any additional text before your response.
		Include any geographic locality where possible.

Current information:
{{.Input}}

Response format:
RESULTS:
Specific Tags: [tag1, tag2, tag3, ...]
Locality: [region1, region2, ...]

Response:`

var sellerPromptTemplate = template.Must(template.New("seller").Parse(sellerPrompt))

var imageSpecialtyPrompt = `First, identify the subject of the image. Then, with the subject of the image in mind, provide a set of 2 - 3 general descriptive tags in CSV format for the professional skills required to create the subject of the image.
You are not describing what's in the image, but rather the type of skills required to produce the subject of the image. Only include tags that you have high confidence about. For instance, if a picture contains high quality cabinetry,
you could emit a tag "custom cabinets". Do not include any additional text before your response. Follow the response format exactly and do not add any additional punctuation. Output CSV only.

Response format:
RESULTS:
Specific Tags: [tag1, tag2, tag3, ...]

Response:
`

var imagePortfolioPrompt = `Describe the high level design properties of this image in terms of specific selling points for the concepts in this image. Give an estimated range of the cost of the work. Give a summary of the level of taste and sophistication.`

type Chat struct {
	model  string
	Client *api.Client
}

func NewChat() *Chat {

	client, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	return &Chat{
		model:  "llama3.2-vision",
		Client: client,
	}
}

// input is all of the messages
func (chat *Chat) BuyerQuery(ctx context.Context, input string) (proposedQuery, followup string) {

	var buf bytes.Buffer
	buyerPromptTemplate.Execute(&buf, struct{ Input string }{input})

	expandedPrompt := buf.String()
	fmt.Println(expandedPrompt)

	req := &api.GenerateRequest{
		Model:  chat.model,
		Prompt: expandedPrompt,

		// set streaming to false
		Stream: new(bool),
	}

	var response string

	respFunc := func(resp api.GenerateResponse) error {
		response = resp.Response

		return nil
	}

	err := chat.Client.Generate(ctx, req, respFunc)

	if err != nil {
		log.Fatal(err)
	}

	return parseBuyerQueryResponse(response)
}

func (chat *Chat) CreatePortfolioDescriptorForImage(ctx context.Context, path string) string {

	imgData, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  chat.model,
		Prompt: imagePortfolioPrompt,
		Images: []api.ImageData{imgData},

		// set streaming to false
		Stream: new(bool),
	}

	var results string

	respFunc := func(resp api.GenerateResponse) error {
		results = resp.Response

		return nil
	}

	err = chat.Client.Generate(ctx, req, respFunc)

	if err != nil {
		log.Fatal(err)
	}

	return results
}

func (chat *Chat) ReturnTagsForImage(ctx context.Context, path string) []string {

	imgData, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	req := &api.GenerateRequest{
		Model:  chat.model,
		Prompt: imageSpecialtyPrompt,
		Images: []api.ImageData{imgData},

		// set streaming to false
		Stream: new(bool),
	}

	return chat.callGenerateAndParseResults(ctx, req)
}

func (chat *Chat) ReturnTagsForString(ctx context.Context, input string) []string {

	var buf bytes.Buffer
	sellerPromptTemplate.Execute(&buf, struct{ Input string }{input})

	expandedPrompt := buf.String()
	fmt.Println(expandedPrompt)

	req := &api.GenerateRequest{
		Model:  chat.model,
		Prompt: expandedPrompt,

		// set streaming to false
		Stream: new(bool),
	}

	return chat.callGenerateAndParseResults(ctx, req)
}

func (chat *Chat) callGenerateAndParseResults(ctx context.Context, req *api.GenerateRequest) []string {
	var results string

	respFunc := func(resp api.GenerateResponse) error {
		results = resp.Response

		return nil
	}

	err := chat.Client.Generate(ctx, req, respFunc)

	if err != nil {
		log.Fatal(err)
	}

	return parseResults(results)
}

func parseBuyerQueryResponse(response string) (proposedQuery, followup string) {

	// Split the response into lines
	lines := strings.Split(strings.TrimSpace(response), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse Query
		if strings.HasPrefix(line, "Search Query:") {
			proposedQuery = strings.TrimSpace(strings.TrimPrefix(line, "Search Query:"))
		}

		// Parse Follow Up Question
		if strings.HasPrefix(line, "Follow Up Question:") {
			followup = strings.TrimSpace(strings.TrimPrefix(line, "Follow Up Question:"))
		}
	}

	return
}
