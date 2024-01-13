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

	// Get supported languages
	{
		languages, err := client.GetLanguages(context.Background(), deepl.Target)
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(languages)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	// Get supported glossary languages pairs
	{
		languages, err := client.GetGlossaryLanguagesPairs(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(languages)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}
