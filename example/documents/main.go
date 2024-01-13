package main

import (
	"context"
	"log"
	"os"

	"github.com/candy12t/go-deepl"
)

func main() {
	authkey := os.Getenv("DEEPL_AUTH_KEY")
	client := deepl.NewClient(authkey)
	ctx := context.Background()

	// upload document
	filePath := "./index.html"
	info, err := client.UploadTranslateDocument(ctx, filePath, "JA", &deepl.TranslateDocumentOption{})
	if err != nil {
		log.Fatal(err)
	}

	// check upload document status
	status, err := client.CheckDocumentStatus(ctx, info.DocumentID, info.DocumentKey)
	if err != nil {
		log.Fatal(err)
	}
	if status.Status != "done" {
		log.Fatalf("not ready for download. translate status is %s.\ndocument_id: %s, document_key: %s", status.Status, info.DocumentID, info.DocumentKey)
	}

	// download translated document
	document, err := client.DownloadTranslatedDocument(ctx, info.DocumentID, info.DocumentKey)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("download.html")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(document); err != nil {
		log.Fatal(err)
	}
}
