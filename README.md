# go-deepl

Unofficial DeepL API client for Go.

## Usage

```go
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
```

## References

- [DeepL API](https://www.deepl.com/en/docs-api)

## APIs

### Translate Text

- [x] POST /v2/translate

### Translate Documents

- [x] POST /v2/document
- [x] POST /v2/document/{document_id}
- [x] POST /v2/document/{document_id}/result

### Manage Glossaries

- [x] GET /v2/glossary-language-pairs
- [x] POST /v2/glossaries
- [x] GET /v2/glossaries
- [x] GET /v2/glossaries/{glossary_id}
- [x] DELETE /v2/glossaries/{glossary_id}
- [x] GET /v2/glossaries/{glossary_id}/entries

### General

- [x] GET /v2/usage
- [x] GET /v2/languages
