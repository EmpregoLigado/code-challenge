package model

import (
	"encoding/json"
	"time"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/pkg/errors"
)

const (
	//StatusActive the identifier for active jobs
	StatusActive string = "active"
	//StatusAny the identifier for any status of a job
	StatusAny string = "any"
	//StatusDraft the identifier for draft jobs
	StatusDraft string = "draft"
	//DateFormat the pattern used to parse the expires date
	DateFormat string = "2/1/2006"
)

//Job represents a job offer
type Job struct {
	PartnerID  int64     `json:"partner_id"`
	CategoryID int64     `json:"category_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
}

//Encrypt a job using the given cipher
func (jb *Job) Encrypt(cipher crypt.Cipher) error {
	crypt, err := cipher.Encrypt(jb.Title)
	if err != nil {
		return errors.Wrap(err, "could not encrypt job")
	}
	jb.Title = crypt
	return nil
}

//Decrypt a job using the given cipher
func (jb *Job) Decrypt(cipher crypt.Cipher) error {
	decrypt, err := cipher.Decrypt(jb.Title)
	if err != nil {
		return errors.Wrap(err, "could not decrypt job")
	}
	jb.Title = decrypt
	return nil
}

//OutputJob represents a job offer with a date in our desired format
type OutputJob struct {
	PartnerID  int64     `json:"partner_id"`
	CategoryID int64     `json:"category_id"`
	Title      string    `json:"title"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
}

//MarshalJSON customizes the date format when outputting json
func (jb *OutputJob) MarshalJSON() ([]byte, error) {
	type Alias Job
	return json.Marshal(&struct {
		ExpiresAt string `json:"expires_at"`
		*Alias
	}{
		ExpiresAt: jb.ExpiresAt.Format(DateFormat),
		Alias:     (*Alias)(jb),
	})
}

//UnmarshalJSON parses dates when unmarshaling json
// func (jb *OutputJob) UnmarshalJSON(data []byte) error {
// 	type Alias Job
// 	var err error
// 	aux := &struct {
// 		ExpiresAt string `json:"expires_at"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(jb),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	jb.ExpiresAt, err = time.Parse(DateFormat, aux.ExpiresAt)
// 	return err
// }
