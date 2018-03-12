package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/EmpregoLigado/code-challenge/model"
	twirp "github.com/twitchtv/twirp"

	pb "github.com/EmpregoLigado/code-challenge/proto"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
)

type jobDataHandler struct {
	storage interfaces.Job
	cipher  crypt.Cipher
}

//NewJobDataHandler creates the abstraction that handles all the job api requests
//All the requests must be encoded as json and encrypted with the pkcs5 algorithm
//The secret key used with the pkcs algorithm must be provided by config or environment
//All the responses will be encoded with the same logic
func NewJobDataHandler(storage interfaces.Job, cipher crypt.Cipher) *jobDataHandler {
	return &jobDataHandler{storage: storage, cipher: cipher}
}

//Common error generator DRY
func (jh *jobDataHandler) requestError(err error) error {
	return twirp.NewError(twirp.InvalidArgument, err.Error())
}

//Create validates and persists a job
func (jh *jobDataHandler) Create(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	jsn, err := jh.cipher.Decrypt(req.Encoded)
	if err != nil {
		return nil, jh.requestError(err)
	}
	job := &model.RequestCreate{}
	err = json.Unmarshal([]byte(jsn), job)
	if err != nil {
		return nil, jh.requestError(err)
	}
	err = jh.storage.Save(&job.Job)
	if err != nil {
		return nil, jh.requestError(err)
	}
	resp := model.ResponseCreate{}
	resp.Status = http.StatusCreated
	encresp, err := crypt.EncryptRequest(jh.cipher, &resp)
	if err != nil {
		return nil, jh.requestError(err)
	}
	return encresp, nil
}

//List retrieves a list of jobs filtering by the Limit,Page and Status params
func (jh *jobDataHandler) List(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	jsn, err := jh.cipher.Decrypt(req.Encoded)
	if err != nil {
		return nil, jh.requestError(err)
	}
	dreq := &model.RequestList{}
	err = json.Unmarshal([]byte(jsn), dreq)
	if err != nil {
		return nil, jh.requestError(err)
	}
	switch dreq.Status {
	case model.StatusActive:
	case model.StatusAny:
	case model.StatusDraft:
	default:
		dreq.Status = model.StatusAny
	}
	if dreq.Limit == 0 {
		dreq.Limit = 100
	}
	lst, err := jh.storage.List(dreq.Limit, dreq.Page, dreq.Status)
	if err != nil {
		return nil, jh.requestError(err)
	}
	res, err := crypt.EncryptRequest(jh.cipher, lst)
	if err != nil {
		return nil, jh.requestError(err)
	}
	return res, nil
}

//Activate changes the status of a job to Active
func (jh *jobDataHandler) Activate(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	jsn, err := jh.cipher.Decrypt(req.Encoded)
	if err != nil {
		return nil, jh.requestError(err)
	}
	dreq := &model.RequestActivate{}
	err = json.Unmarshal([]byte(jsn), dreq)
	if err != nil {
		return nil, jh.requestError(err)
	}
	err = jh.storage.Activate(dreq.PartnerID)
	if err != nil {
		return nil, jh.requestError(err)
	}
	resp := model.ResponseActivate{}
	resp.Status = http.StatusOK
	encresp, err := crypt.EncryptRequest(jh.cipher, &resp)
	if err != nil {
		return nil, jh.requestError(err)
	}
	return encresp, nil
}

//Percentage totalizes the percentage of active jobs of each category
func (jh *jobDataHandler) Percentage(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	jsn, err := jh.cipher.Decrypt(req.Encoded)
	if err != nil {
		return nil, jh.requestError(err)
	}
	dreq := &model.RequestPercentage{}
	err = json.Unmarshal([]byte(jsn), dreq)
	if err != nil {
		return nil, jh.requestError(err)
	}
	perc, err := jh.storage.Percentage(dreq.CategoryID)
	if err != nil {
		return nil, jh.requestError(err)
	}
	res, err := crypt.EncryptRequest(jh.cipher, perc)
	if err != nil {
		return nil, jh.requestError(err)
	}
	return res, nil
}
