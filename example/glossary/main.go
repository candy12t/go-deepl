package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/candy12t/go-deepl"
)

func main() {
	authkey := os.Getenv("DEEPL_AUTH_KEY")
	client := deepl.NewClient(authkey)
	ctx := context.Background()

	// Create Glossary
	{
		entries := deepl.NewGlossaryEntries(
			deepl.GlossaryEntry{Source: "Hello world", Target: "こんにちは、世界"},
		)
		glossary, err := client.CreateGlossary(ctx, "Hello world", "EN", "JA", entries.EncodeGlossaryEntries(deepl.TSV), deepl.TSV)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(glossary)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	// Get All Glossaries
	{
		glossaries, err := client.GetGlossaries(ctx)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(glossaries)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	// Get Glossary of specified id
	{
		id := "glossary_id"
		glossary, err := client.GetGlossary(ctx, id)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(glossary)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	// Delete Glossary of specified id
	{
		id := "glossary_id"
		err := client.DeleteGlossary(ctx, id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
