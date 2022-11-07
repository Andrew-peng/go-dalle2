package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Andrew-peng/go-dalle2/dalle2"
)

func main() {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		log.Fatal("Environment variable OPENAI_API_KEY is not set")
	}
	client, err := dalle2.MakeNewClientV1(apiKey)
	if err != nil {
		log.Fatalf("Error initializing client: %s", err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	imgPath := filepath.Join(curDir, "edit/otter.png")
	maskPath := filepath.Join(curDir, "edit/mask.png")
	resp, err := client.Edit(
		context.Background(),
		imgPath,
		maskPath,
		"A cute baby sea otter wearing a large sombrero",
		dalle2.WithNumImages(1),
		dalle2.WithSize(dalle2.SMALL),
		dalle2.WithFormat(dalle2.URL),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created: %d\n", resp.Created)
	fmt.Println("Images:")
	for _, img := range resp.Data {
		fmt.Printf("\t%s\n", img.Url)
	}
}
