package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// AppConfig holds all configuration values.
type AppConfig struct {
	Vapi      VapiConfig
	Assistant AssistantConfig
}

type VapiConfig struct {
	APIKey  string
	BaseURL string
}

type AssistantConfig struct {
	Name             string
	FirstMessage     string
	FirstMessageMode string
	Model            ModelConfig
	Voice            VoiceConfig
	Transcriber      TranscriberConfig
}

type ModelConfig struct {
	Provider     string
	Model        string
	Temperature  float64
	SystemPrompt string
}

type VoiceConfig struct {
	Provider string
	VoiceID  string
}

type TranscriberConfig struct {
	Provider string
	Model    string
	Language string
}

// Load reads config.yaml and environment variables, then returns AppConfig.
func Load() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..") // cmd/ dizininden çalıştırılırsa

	// Environment variable override
	_ = viper.BindEnv("vapi.api_key", "VAPI_API_KEY")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config dosyası okunamadı: %v", err)
	}

	fmt.Printf("Config yüklendi: %s\n", viper.ConfigFileUsed())

	return &AppConfig{
		Vapi: VapiConfig{
			APIKey:  viper.GetString("vapi.api_key"),
			BaseURL: viper.GetString("vapi.base_url"),
		},
		Assistant: AssistantConfig{
			Name:             viper.GetString("assistant.name"),
			FirstMessage:     viper.GetString("assistant.first_message"),
			FirstMessageMode: viper.GetString("assistant.first_message_mode"),
			Model: ModelConfig{
				Provider:     viper.GetString("assistant.model.provider"),
				Model:        viper.GetString("assistant.model.model"),
				Temperature:  viper.GetFloat64("assistant.model.temperature"),
				SystemPrompt: viper.GetString("assistant.model.system_prompt"),
			},
			Voice: VoiceConfig{
				Provider: viper.GetString("assistant.voice.provider"),
				VoiceID:  viper.GetString("assistant.voice.voice_id"),
			},
			Transcriber: TranscriberConfig{
				Provider: viper.GetString("assistant.transcriber.provider"),
				Model:    viper.GetString("assistant.transcriber.model"),
				Language: viper.GetString("assistant.transcriber.language"),
			},
		},
	}
}
