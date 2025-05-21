package ollama

import (
	"context"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
	"os"
)

func AskOllama(prompt string, u *[]api.Message) string {

	c := api.NewClient(
		&url.URL{Scheme: "http", Host: os.Getenv("OLLAMA_HOST")},
		http.DefaultClient,
	)
	stream := false
	sysprompt, err := os.ReadFile("sysprompt.txt")
	if err != nil {
		panic(err)
	}
	var result string
	*u = append(*u, api.Message{
		Role:    "system",
		Content: string(sysprompt),
	})

	*u = append(*u, api.Message{
		Role:    "user",
		Content: prompt,
	})
	e := c.Chat(
		context.Background(),
		&api.ChatRequest{
			Model:    "llama3.1",
			Messages: *u,
			Stream:   &stream,
		},
		func(response api.ChatResponse) error {
			*u = append(*u, api.Message{
				Role:    "assistant",
				Content: response.Message.Content,
			})
			result += response.Message.Content // accumulate in case of streaming or chunked response
			return nil
		},
	)
	if e != nil {
		panic(e)
	}
	return result
}
