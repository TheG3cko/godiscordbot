package main

import (
	"context"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
	"os"
)

func askollamaNoHistory(prompt string) string {
	c := api.NewClient(
		&url.URL{Scheme: "http", Host: "localhost:11434"},
		http.DefaultClient,
	)
	stream := false
	var result string
	e := c.Generate(
		context.Background(),
		&api.GenerateRequest{
			Model:  "llama3.2",
			Prompt: prompt,
			Stream: &stream,
		},
		func(response api.GenerateResponse) error {
			result += response.Response // accumulate in case of streaming or chunked response
			return nil
		},
	)

	if e != nil {
		panic(e)
	}
	return result
}

func askollama(prompt string) string {

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
	userHistories = append(userHistories, api.Message{
		Role:    "system",
		Content: string(sysprompt),
	})

	userHistories = append(userHistories, api.Message{
		Role:    "user",
		Content: prompt,
	})
	e := c.Chat(
		context.Background(),
		&api.ChatRequest{
			Model:    "llama3.2",
			Messages: userHistories,
			Stream:   &stream,
		},
		func(response api.ChatResponse) error {
			userHistories = append(userHistories, api.Message{
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
