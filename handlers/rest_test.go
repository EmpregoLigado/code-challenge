package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/EmpregoLigado/code-challenge/crypt/null"
	"github.com/EmpregoLigado/code-challenge/model"
	pb "github.com/EmpregoLigado/code-challenge/proto"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
	"github.com/EmpregoLigado/code-challenge/storage/memory"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	dataserver, restserver *httptest.Server
	storage                interfaces.Job
	bearer                 string
	joblist                []model.Job
)

const (
	secret string = "aeae42cd8f444313a4f300088713e71c"
)

func TestMain(m *testing.M) {
	//Initializes a dummy data server
	cipher := null.Cipher{}
	storage = memory.NewJob(cipher)
	svcHandler := pb.NewJobServer(NewJobDataHandler(storage, cipher), nil)
	smux := http.NewServeMux()
	smux.Handle(pb.JobPathPrefix, svcHandler)
	dataserver := httptest.NewServer(smux)

	//Inicializes the testing rest server
	backend := pb.NewJobProtobufClient(dataserver.URL, &http.Client{})

	twentyFourHours := time.Now().UTC().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": twentyFourHours.Unix(),
	})
	bearer, _ = token.SignedString([]byte(secret))

	//Our list of fake jobs to test
	joblist = []model.Job{
		model.Job{
			PartnerID:  1009,
			CategoryID: 9001,
			Title:      "Fake Job One",
			ExpiresAt:  twentyFourHours,
		},
		model.Job{
			PartnerID:  1010,
			CategoryID: 9001,
			Title:      "Fake Job Two",
			ExpiresAt:  twentyFourHours.Add(1 * time.Hour),
		},
		model.Job{
			PartnerID:  1011,
			CategoryID: 9002,
			Title:      "Fake Job Three",
			ExpiresAt:  twentyFourHours.Add(2 * time.Hour),
		},
		model.Job{
			PartnerID:  1012,
			CategoryID: 9003,
			Title:      "Fake Job Four",
			ExpiresAt:  twentyFourHours.Add(3 * time.Hour),
		},
		model.Job{
			PartnerID:  1013,
			CategoryID: 9003,
			Title:      "Fake Job Five",
			ExpiresAt:  twentyFourHours.Add(4 * time.Hour),
		},
	}

	mux := NewRestHandler(cipher, backend, secret)
	restserver = httptest.NewServer(mux)
	retCode := m.Run()

	os.Exit(retCode)
}

func newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not create request for %s", err.Error()))
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+bearer)
	return req, nil
}

func TestCreatJobInputValidation(t *testing.T) {
	surl := restserver.URL + "/jobs"
	form := make(url.Values)

	expected := func(surl string, form url.Values, code int, message string) {
		req, err := newRequest(http.MethodPost, surl, strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

		client := http.Client{}
		r, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), message) {
			t.Fatalf("expecting error \"%s\", got \"%s\"", message, body)
		}
		if r.StatusCode != code {
			t.Fatalf("expecting status %d, got %d, %s", code, r.StatusCode, body)
		}
	}

	expected(surl, form, http.StatusBadRequest, "missing partner_id value")

	form["partner_id"] = append(form["partner_id"], "AB")
	expected(surl, form, http.StatusBadRequest, "invalid partner_id value")

	form["partner_id"][0] = "1009"
	expected(surl, form, http.StatusBadRequest, "missing category_id value")

	form["category_id"] = append(form["category_id"], "AB")
	expected(surl, form, http.StatusBadRequest, "invalid category_id value")

	form["category_id"][0] = "901"
	expected(surl, form, http.StatusBadRequest, "missing or empty title given")

	form["title"] = append(form["title"], "Job Title")
	expected(surl, form, http.StatusBadRequest, "missing expires_at value")

	form["expires_at"] = append(form["expires_at"], "AB")
	expected(surl, form, http.StatusBadRequest, "invalid expiration date")

	form["expires_at"][0] = "1/1/1990"
	expected(surl, form, http.StatusBadRequest, "job already expired")

	form["expires_at"][0] = time.Now().Format(model.DateFormat)
	expected(surl, form, http.StatusCreated, "")

	err := storage.Delete(1009)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreatJobs(t *testing.T) {
	surl := restserver.URL + "/jobs"
	createJob := func(jb *model.Job) (*http.Response, error) {
		form := make(url.Values)

		form["partner_id"] = append(form["partner_id"], fmt.Sprintf("%d", jb.PartnerID))
		form["category_id"] = append(form["category_id"], fmt.Sprintf("%d", jb.CategoryID))
		form["title"] = append(form["title"], jb.Title)
		form["expires_at"] = append(form["expires_at"], jb.ExpiresAt.Format(model.DateFormat))

		req, err := newRequest(http.MethodPost, surl, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

		client := http.Client{}
		return client.Do(req)
	}

	for _, job := range joblist {
		resp, err := createJob(&job)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expecting status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	}

	//Since we dont have a GET method on the API, test all the cases against the storage
	//Thanks dependency injection
	for _, job := range joblist {
		stjob, err := storage.Get(job.PartnerID)
		if err != nil {
			t.Fatal(err)
		}
		if stjob.CategoryID != job.CategoryID {
			t.Fatalf("job saved with wrong category id, expected %d, got %d", job.CategoryID, stjob.CategoryID)
		}
		if stjob.Status != model.StatusDraft {
			t.Fatalf("job saved with wrong status, expected %s, got %s", model.StatusDraft, stjob.Status)
		}
		if !stjob.ExpiresAt.Equal(stjob.ExpiresAt) {
			t.Fatalf("job saved with wrong expiration date, expected %s, got %s",
				job.ExpiresAt.Format(model.DateFormat),
				stjob.ExpiresAt.Format(model.DateFormat))
		}
	}
}

func TestListJobs(t *testing.T) {
	//To parse the rest service output format we need the expires date to be a string
	//Since we are not using this anywere outside this test, declare it locally
	type OutJob struct {
		PartnerID  int64  `json:"partner_id"`
		CategoryID int64  `json:"category_id"`
		Title      string `json:"title"`
		Status     string `json:"status"`
		ExpiresAt  string `json:"expires_at"`
	}

	req, err := newRequest(http.MethodGet, restserver.URL+"/jobs", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	jobs := []OutJob{}
	err = decoder.Decode(&jobs)
	if err != nil {
		t.Fatal(err)
	}

	if len(jobs) < len(joblist) {
		t.Fatalf("expecting %d jobs, got %d", len(joblist), len(jobs))
	}

	for i, job := range jobs {
		if job.Status != model.StatusDraft {
			t.Fatal("new jobs should always have a draft status")
		}
		if joblist[i].CategoryID != job.CategoryID {
			t.Fatalf("job returned with wrong category, expected %d, got %d", joblist[i].CategoryID, job.CategoryID)
		}
		if joblist[i].ExpiresAt.Format(model.DateFormat) != job.ExpiresAt {
			t.Fatalf("job returned with wrong expiration date, expected %s, got %s",
				joblist[i].ExpiresAt.Format(model.DateFormat),
				job.ExpiresAt)
		}
		if joblist[i].Title != job.Title {
			t.Fatalf("job returned with wrong title, expected %s, got %s",
				joblist[i].Title, job.Title)
		}
	}
}

func TestActivateJob(t *testing.T) {
	activateJob := func(id int64) (*http.Response, error) {
		surl := fmt.Sprintf("%s/jobs/%d/activate", restserver.URL, id)
		req, err := newRequest(http.MethodPost, surl, nil)
		if err != nil {
			t.Fatal(err)
		}
		client := http.Client{}
		return client.Do(req)
	}

	resp, err := activateJob(joblist[0].PartnerID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = activateJob(6667)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expecting status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

}
func TestPercentage(t *testing.T) {
	req, err := newRequest(http.MethodGet, restserver.URL+"/category/9001", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Del("Authorization")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	perc := model.CategoryPercentage{}
	err = decoder.Decode(&perc)
	if err != nil {
		t.Fatal(err)
	}
	if perc.Available != 1 {
		t.Fatalf("should have only one job active in category 9001, got %d", perc.Available)
	}
	if perc.Percentage != 50.00 {
		t.Fatalf("should have 50%% of jobs active, got %f", perc.Percentage)
	}

	req, err = newRequest(http.MethodGet, restserver.URL+"/category/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Del("Authorization")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()
	decoder = json.NewDecoder(resp.Body)
	perc = model.CategoryPercentage{}
	err = decoder.Decode(&perc)
	if err != nil {
		t.Fatal(err)
	}
	if perc.Available != 0 {
		t.Fatal("should return zero jobs")
	}
	if perc.Percentage != 0.00 {
		t.Fatalf("should return zero on empty category list")
	}

}
