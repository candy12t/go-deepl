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

	usage, err := client.GetUsage(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(usage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
