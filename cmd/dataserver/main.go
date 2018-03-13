package main

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/EmpregoLigado/code-challenge/crypt/pkcs5"
	"github.com/EmpregoLigado/code-challenge/handlers"
	"github.com/EmpregoLigado/code-challenge/model"
	pb "github.com/EmpregoLigado/code-challenge/proto"
	"github.com/EmpregoLigado/code-challenge/storage"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
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
		port = "8079"
	}
	//TODO: Minimun sanity check on port value

	db := os.Getenv("DB")
	if db == "" {
		db = "mysql://root:codechallenge@/jobdb"
	}

	bootstrap := os.Getenv("BOOTSTRAP")
	if bootstrap == "" {
		bootstrap = "jobs.txt"
	}

	cipher := pkcs5.Cipher{Key: bkey}
	storage, err := storage.NewJob(cipher, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if _, err := os.Stat(bootstrap); err == nil {
		go loadFile(bootstrap, storage)
	}

	svcHandler := pb.NewJobServer(handlers.NewJobDataHandler(storage, cipher), nil)

	// You can use any mux you like - JobServer gives you an http.Handler.
	mux := http.NewServeMux()
	// The generated code includes a const, <ServiceName>PathPrefix, which
	// can be used to mount your service on a mux.
	mux.Handle(pb.JobPathPrefix, svcHandler)
	http.ListenAndServe(":"+port, mux)
}

func loadFile(path string, storage interfaces.Job) {
	fmt.Printf("\nImporting file %s\n", path)
	in, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(in)
	r.Comma = '|'

	now := time.Now()
	//skips header
	r.Read()
	var i, j int64

	for {
		i++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) < 4 {
			fmt.Fprintf(os.Stderr, "%d: expecting 4 fields, got %d\n", i, len(record))
			continue
		}
		job := &model.Job{}
		job.PartnerID, err = strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d: invalid partnerId value %s, %s\n", i, record[0], err.Error())
			continue
		}
		job.Title = record[1]
		if job.Title == "" {
			fmt.Fprintf(os.Stderr, "%d: missing or empty title given\n", i)
			continue
		}
		job.CategoryID, err = strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d: invalid categoryId value %s, %s\n", i, record[2], err.Error())
			continue
		}
		job.ExpiresAt, err = time.ParseInLocation(model.DateFormat, record[3], now.Location())
		if err != nil || job.ExpiresAt.IsZero() {
			fmt.Fprintf(os.Stderr, "%d: invalid expiration date\n", i)
			continue
		}
		err = storage.Save(job)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%d: error saving job, %s\n", i, err.Error())
		}
		j++
	}
	fmt.Printf("Import proccess finished with %d valid records of %d found\n", j, i)
}
