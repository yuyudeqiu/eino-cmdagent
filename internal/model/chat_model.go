package model

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
)

func NewChatModel(ctx context.Context) (*ollama.ChatModel, error) {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen2.5:7b",
		// Model: "qwen2.5:14b",
	})
	if err != nil {
		return nil, err
	}
	return chatModel, nil
}

func NewArkModel(ctx context.Context) (*ark.ChatModel, error) {
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("ARK_MODEL"),
	})
	if err != nil {
		return nil, err
	}
	return model, nil
}
