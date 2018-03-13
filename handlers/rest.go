package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/EmpregoLigado/code-challenge/middleware"
	"github.com/EmpregoLigado/code-challenge/model"
	"github.com/EmpregoLigado/code-challenge/proto"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
)

//RestHandler handles all the REST requests for job manipulation
type RestHandler struct {
	cipher  crypt.Cipher
	backend proto.Job
}

//NewRestHandler initializes a new mux for the given handler
func NewRestHandler(cipher crypt.Cipher, backend proto.Job, secret string) http.Handler {
	handler := RestHandler{cipher, backend}

	mux := mux.NewRouter()
	jwtsec := middleware.JWTSecure(secret)
	mux.Path("/jobs").Methods(http.MethodPost).Handler(jwtsec(handler.Create))
	mux.Path("/jobs").Methods(http.MethodGet).Handler(jwtsec(handler.List))
	mux.Path("/jobs/{jobid}/activate").Methods(http.MethodPost).Handler(jwtsec(handler.Activate))
	mux.Path("/category/{categoryid}").Methods(http.MethodGet).Handler(http.HandlerFunc(handler.Percentage))
	return mux
}

//Create receive, validate and persists a job.
//It espects a regular form post with the fields:
// partner_id: int
// category_id: int
// title: string
// expires_at: date DD/MM/YYYY
func (rh *RestHandler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeError(w, http.StatusBadRequest, "unable to parse form values")
		return
	}

	job := model.Job{}
	if r.FormValue("partner_id") == "" {
		writeError(w, http.StatusBadRequest, "missing partner_id value")
		return
	}
	job.PartnerID, err = strconv.ParseInt(r.FormValue("partner_id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid partner_id value")
		return
	}
	if r.FormValue("category_id") == "" {
		writeError(w, http.StatusBadRequest, "missing category_id value")
		return
	}
	job.CategoryID, err = strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid category_id value").Error())
		return
	}
	job.Title = strings.TrimSpace(r.FormValue("title"))
	if job.Title == "" {
		writeError(w, http.StatusBadRequest, "missing or empty title given")
		return
	}

	if r.FormValue("expires_at") == "" {
		writeError(w, http.StatusBadRequest, "missing expires_at value")
		return
	}
	now := time.Now()
	job.ExpiresAt, err = time.ParseInLocation(model.DateFormat, r.FormValue("expires_at"), now.Location())
	if err != nil || job.ExpiresAt.IsZero() {
		writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid expiration date").Error())
		return
	}
	//Times are parsed without hour, so whe have to add the hours until de end of the day
	job.ExpiresAt = job.ExpiresAt.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	if job.ExpiresAt.Before(now) {
		writeError(w, http.StatusBadRequest, "job already expired")
		return
	}
	req := model.RequestCreate{}
	req.Job = job

	encreq, err := crypt.EncryptRequest(rh.cipher, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = rh.backend.Create(context.Background(), encreq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, nil)
}

func (rh *RestHandler) List(w http.ResponseWriter, r *http.Request) {
	req := model.RequestList{}
	var err error
	vars := r.URL.Query()

	limit, ok := vars["limit"]
	if !ok {
		req.Limit = 100
	} else {
		req.Limit, err = strconv.Atoi(limit[0])
		if err != nil || req.Limit < 0 {
			writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid limit value").Error())
			return
		}
	}

	page, ok := vars["page"]
	if !ok {
		req.Page = 0
	} else {
		req.Page, err = strconv.Atoi(page[0])
		if err != nil || req.Page < 0 {
			writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid page value").Error())
			return
		}
	}

	req.Status = model.StatusAny
	status, ok := vars["status"]
	if ok {
		switch strings.ToLower(status[0]) {
		case model.StatusActive:
			req.Status = model.StatusActive
		case model.StatusDraft:
			req.Status = model.StatusDraft
		}
	}

	encreq, err := crypt.EncryptRequest(rh.cipher, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := rh.backend.List(context.Background(), encreq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsn, err := rh.cipher.Decrypt(res.Encoded)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resplist := []model.OutputJob{}
	err = json.Unmarshal([]byte(jsn), &resplist)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, 200, resplist)
}

func (rh *RestHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id, ok := vars["jobid"]
	if !ok {
		writeError(w, http.StatusBadRequest, "missing partner_id value")
		return
	}
	req := model.RequestActivate{}
	req.PartnerID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid partner_id value").Error())
		return
	}

	encreq, err := crypt.EncryptRequest(rh.cipher, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = rh.backend.Activate(context.Background(), encreq)
	if err != nil {
		if strings.Contains(err.Error(), interfaces.ErrJobNotFound.Error()) {
			writeError(w, http.StatusNotFound, interfaces.ErrJobNotFound.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (rh *RestHandler) Percentage(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id, ok := vars["categoryid"]
	if !ok {
		writeError(w, http.StatusBadRequest, "missing category_id value")
		return
	}
	req := model.RequestPercentage{}
	req.CategoryID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.Wrap(err, "invalid category_id value").Error())
		return
	}

	encreq, err := crypt.EncryptRequest(rh.cipher, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := rh.backend.Percentage(context.Background(), encreq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsn, err := rh.cipher.Decrypt(res.Encoded)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	catperc := &model.CategoryPercentage{}
	err = json.Unmarshal([]byte(jsn), catperc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, 200, catperc)
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	var msgerr = struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(msgerr)
}

func writeResponse(w http.ResponseWriter, code int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(message)
}
