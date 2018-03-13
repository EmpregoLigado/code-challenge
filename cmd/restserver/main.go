package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/EmpregoLigado/code-challenge/crypt/pkcs5"
	"github.com/EmpregoLigado/code-challenge/handlers"
	"github.com/EmpregoLigado/code-challenge/proto"
)

func main() {
	key := os.Getenv("KEY")
	if key == "" {
		key = "01020304050607080910111213141516"
	}
	bkey, err := hex.DecodeString(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" is not a valid hex key, %s", key, err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//TODO: Minimun sanity check on port value

	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "ffeaa4bbe7d741a9cc900ddabb7961a2d44a9f10"
	}

	datahost := os.Getenv("DATA_HOST")
	if datahost == "" {
		datahost = "http://data:8080"
	}

	cipher := pkcs5.Cipher{Key: bkey}
	backend := proto.NewJobProtobufClient(datahost, &http.Client{})
	mux := handlers.NewRestHandler(cipher, backend, secret)
	http.ListenAndServe(":"+port, mux)
}
