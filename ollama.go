package main

import (
	"context"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
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

type Message struct {
	Content string
	Role    string
}

func askollama(prompt string) string {
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

	c.Chat(
		context.Background(),
		&api.ChatRequest{
			Model: "llama3.2",
			Messages: []api.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is the capital of France?"},
			},
			Stream:    nil,
			Format:    nil,
			KeepAlive: nil,
			Tools:     nil,
			Options:   nil,
		},
		nil,
	)

	if e != nil {
		panic(e)
	}
	return result
}
