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

func askollama(prompt string) string {

	c := api.NewClient(
		&url.URL{Scheme: "http", Host: "10.10.10.24:11434"},
		http.DefaultClient,
	)
	stream := false
	var result string
	userHistories = append(userHistories, api.Message{
		Role:    "system",
		Content: "You are a Discord Chatbot. At the beginning of each message, there is the username of the user. This way you can separate them. ",
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
