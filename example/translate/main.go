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

	translatetext, err := client.TranslateText(context.Background(), []string{"Hello world"}, "JA", deepl.TranslateOption{SourceLang: "EN"})
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(translatetext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
