package crypt

import (
	"encoding/json"

	pb "github.com/EmpregoLigado/code-challenge/proto"
	"github.com/pkg/errors"
)

//Cipher implements the required functions to encrypt data on various parts of the project
type Cipher interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

//EncryptRequest request encoder for IPC
func EncryptRequest(cipher Cipher, req interface{}) (*pb.Payload, error) {
	res := &pb.Payload{}
	json, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal request")
	}
	res.Encoded, err = cipher.Encrypt(string(json))
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal request")
	}
	return res, nil
}
