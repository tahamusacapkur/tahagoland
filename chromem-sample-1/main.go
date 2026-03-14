/*
ICI Tech Teknoloji A.Ş.
User        : ICI
Name        : Ibrahim COBANI
Date        : 7.03.2026
Time        : 15:31
Notes       :
*/
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/philippgille/chromem-go"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	// Persistent database - saves to disk
	db, err := chromem.NewPersistentDB("./chromem_data", false)
	if err != nil {
		panic(err)
	}

	// Try to get existing collection, if not exists create new one
	c := db.GetCollection("knowledge-base", nil)
	if c == nil {
		c, err = db.CreateCollection("knowledge-base", nil, nil)
		if err != nil {
			panic(err)
		}

		err = c.AddDocuments(ctx, []chromem.Document{
			{
				ID:      "1",
				Content: "The sky is blue because of Rayleigh scattering.",
			},
			{
				ID:      "2",
				Content: "Leaves are green because chlorophyll absorbs red and blue light.",
			},
		}, runtime.NumCPU())
		if err != nil {
			panic(err)
		}
		fmt.Println("Yeni koleksiyon oluşturuldu ve dökümanlar eklendi.")
	} else {
		fmt.Println("Mevcut koleksiyon diskten yüklendi.")
	}

	fmt.Println("\nArama CLI'ya hoş geldiniz! (Çıkmak için Ctrl+C)")
	fmt.Println(strings.Repeat("-", 50))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nSorgunuz: ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}

		startTime := time.Now()
		res, err := c.Query(ctx, query, 1, nil, nil)
		if err != nil {
			fmt.Printf("Hata: %v\n", err)
			continue
		}

		elapsed := time.Since(startTime)
		fmt.Printf("ID: %v\nBenzerlik: %v\nİçerik: %v\nSüre: %v\n", res[0].ID, res[0].Similarity, res[0].Content, elapsed)
	}
}
