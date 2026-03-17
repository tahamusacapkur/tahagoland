package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"vapi-example-1/internal/client"
	"vapi-example-1/internal/config"
	"vapi-example-1/internal/model"
)

func main() {
	cfg := config.Load()

	if cfg.Vapi.APIKey == "" || cfg.Vapi.APIKey == "YOUR_VAPI_API_KEY_HERE" {
		fmt.Println("Lütfen config.yaml dosyasında vapi.api_key değerini güncelleyin!")
		fmt.Println("Veya VAPI_API_KEY environment variable'ı set edin.")
		fmt.Println()
		fmt.Println("API Key almak için: https://dashboard.vapi.ai → Organization Settings → API Keys")
		os.Exit(1)
	}

	vc := client.New(cfg.Vapi.APIKey, cfg.Vapi.BaseURL)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "create":
		runCreate(vc, cfg)
	case "list":
		runList(vc)
	case "get":
		if len(os.Args) < 3 {
			log.Fatal("Kullanım: go run ./cmd create | list | get <id> | delete <id>")
		}
		runGet(vc, os.Args[2])
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Kullanım: go run ./cmd delete <assistant_id>")
		}
		runDelete(vc, os.Args[2])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Vapi Voice AI - Go Demo")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("Kullanım:")
	fmt.Println("  go run ./cmd create          Yeni asistan oluştur")
	fmt.Println("  go run ./cmd list             Tüm asistanları listele")
	fmt.Println("  go run ./cmd get <id>         Asistan detayını getir")
	fmt.Println("  go run ./cmd delete <id>      Asistanı sil")
}

func runCreate(vc *client.VapiClient, cfg *config.AppConfig) {
	req := model.CreateAssistantRequest{
		Name:             cfg.Assistant.Name,
		FirstMessage:     cfg.Assistant.FirstMessage,
		FirstMessageMode: cfg.Assistant.FirstMessageMode,
		Model: model.Model{
			Provider:      cfg.Assistant.Model.Provider,
			Model:         cfg.Assistant.Model.Model,
			Temperature:   cfg.Assistant.Model.Temperature,
			SystemMessage: cfg.Assistant.Model.SystemPrompt,
		},
		Voice: model.Voice{
			Provider: cfg.Assistant.Voice.Provider,
			VoiceID:  cfg.Assistant.Voice.VoiceID,
		},
		Transcriber: &model.Transcriber{
			Provider: cfg.Assistant.Transcriber.Provider,
			Model:    cfg.Assistant.Transcriber.Model,
			Language: cfg.Assistant.Transcriber.Language,
		},
	}

	fmt.Printf("Asistan oluşturuluyor: %s\n", req.Name)

	assistant, err := vc.CreateAssistant(req)
	if err != nil {
		log.Fatalf("Asistan oluşturulamadı: %v", err)
	}

	fmt.Println("Asistan başarıyla oluşturuldu!")
	fmt.Printf("  ID:      %s\n", assistant.ID)
	fmt.Printf("  Name:    %s\n", assistant.Name)
	fmt.Printf("  Created: %s\n", assistant.CreatedAt.Format(time.RFC3339))
}

func runList(vc *client.VapiClient) {
	assistants, err := vc.ListAssistants()
	if err != nil {
		log.Fatalf("Asistanlar listelenemedi: %v", err)
	}

	if len(assistants) == 0 {
		fmt.Println("Henüz hiç asistan yok.")
		return
	}

	fmt.Printf("Toplam %d asistan bulundu:\n\n", len(assistants))
	for i, a := range assistants {
		fmt.Printf("  %d. [%s] %s\n", i+1, a.ID, a.Name)
	}
}

func runGet(vc *client.VapiClient, id string) {
	assistant, err := vc.GetAssistant(id)
	if err != nil {
		log.Fatalf("Asistan bulunamadı: %v", err)
	}

	fmt.Println("Asistan Detayı:")
	fmt.Printf("  ID:            %s\n", assistant.ID)
	fmt.Printf("  Name:          %s\n", assistant.Name)
	fmt.Printf("  First Message: %s\n", assistant.FirstMessage)
	fmt.Printf("  Created:       %s\n", assistant.CreatedAt.Format(time.RFC3339))
}

func runDelete(vc *client.VapiClient, id string) {
	if err := vc.DeleteAssistant(id); err != nil {
		log.Fatalf("Asistan silinemedi: %v", err)
	}

	fmt.Printf("Asistan silindi: %s\n", id)
}
