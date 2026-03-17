package model

import "time"

// Assistant represents a Vapi voice assistant.
type Assistant struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	FirstMessage     string    `json:"firstMessage"`
	FirstMessageMode string    `json:"firstMessageMode"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Model represents the LLM configuration for an assistant.
type Model struct {
	Provider      string  `json:"provider"`
	Model         string  `json:"model"`
	Temperature   float64 `json:"temperature,omitempty"`
	SystemMessage string  `json:"systemMessage,omitempty"`
}

// Voice represents the text-to-speech configuration.
type Voice struct {
	Provider string `json:"provider"`
	VoiceID  string `json:"voiceId"`
}

// Transcriber represents the speech-to-text configuration.
type Transcriber struct {
	Provider string `json:"provider"`
	Model    string `json:"model,omitempty"`
	Language string `json:"language,omitempty"`
}

// CreateAssistantRequest is the payload for creating a new assistant.
type CreateAssistantRequest struct {
	Name             string       `json:"name"`
	FirstMessage     string       `json:"firstMessage,omitempty"`
	FirstMessageMode string       `json:"firstMessageMode,omitempty"`
	Model            Model        `json:"model"`
	Voice            Voice        `json:"voice"`
	Transcriber      *Transcriber `json:"transcriber,omitempty"`
}
