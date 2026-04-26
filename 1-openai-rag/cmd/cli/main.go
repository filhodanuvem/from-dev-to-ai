package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: cli <path-to-pdf>")
	}

	pdfPath := os.Args[1]
	apiKey := os.Getenv("OPENAI_KEY")
	vectorStoreID := os.Getenv("OPENAI_VECTOR_DB_ID")

	if apiKey == "" || vectorStoreID == "" {
		log.Fatal("OPENAI_KEY and OPENAI_VECTOR_DB_ID must be set")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	f, err := os.Open(pdfPath)
	if err != nil {
		log.Fatalf("open pdf: %v", err)
	}
	defer f.Close()

	fmt.Printf("Uploading %s...\n", pdfPath)
	uploadedFile, err := client.Files.New(ctx, openai.FileNewParams{
		File:    f,
		Purpose: openai.FilePurposeAssistants,
	})
	if err != nil {
		log.Fatalf("upload file: %v", err)
	}
	fmt.Printf("File uploaded: %s\n", uploadedFile.ID)

	fmt.Printf("Attaching to vector store %s...\n", vectorStoreID)
	_, err = client.VectorStores.Files.New(ctx, vectorStoreID, openai.VectorStoreFileNewParams{
		FileID: uploadedFile.ID,
	})
	if err != nil {
		log.Fatalf("attach to vector store: %v", err)
	}

	for {
		vsFile, err := client.VectorStores.Files.Get(ctx, vectorStoreID, uploadedFile.ID)
		if err != nil {
			log.Fatalf("poll status: %v", err)
		}
		fmt.Printf("Status: %s\n", vsFile.Status)
		switch vsFile.Status {
		case openai.VectorStoreFileStatusCompleted:
			fmt.Println("Done! File is indexed in the vector store.")
			return
		case openai.VectorStoreFileStatusInProgress:
			time.Sleep(time.Second)
		default:
			log.Fatalf("unexpected status: %s", vsFile.Status)
		}
	}
}
