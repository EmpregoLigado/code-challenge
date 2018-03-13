package sql

import (
	"database/sql"
	"fmt"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/EmpregoLigado/code-challenge/model"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
	"github.com/pkg/errors"
)

//Job base type to access in memory jobs
type Job struct {
	db     *sql.DB
	cipher crypt.Cipher
}

//NewJob creates a new instance of a memory job storage
func NewJob(cipher crypt.Cipher, db *sql.DB) *Job {
	job := &Job{}
	job.cipher = cipher
	job.db = db
	return job
}

//Get retrieves one job by partnerID
func (jb *Job) Get(ID int64) (*model.Job, error) {
	var res model.Job
	row := jb.db.QueryRow("SELECT partner_id,title,category_id,expires_at,status FROM jobs WHERE partner_id=?", ID)
	err := row.Scan(&res.PartnerID, &res.Title, &res.CategoryID, &res.ExpiresAt, &res.Status)
	if err == sql.ErrNoRows {
		return nil, interfaces.ErrJobNotFound
	}
	return &res, err
}

//Save persists in memory one job
func (jb *Job) Save(job *model.Job) error {
	if job.PartnerID == 0 {
		return errors.New("missing partner ID on job creation")
	}
	if job.CategoryID == 0 {
		return errors.New("missing category ID on job creation")
	}
	//Enforces status based on the previous value, otherwise it's a draft
	prev, _ := jb.Get(job.PartnerID)

	var query string
	var err error
	job.Status = model.StatusDraft
	err = job.Encrypt(jb.cipher)
	if err != nil {
		return errors.Wrap(err, "error saving job")
	}

	if prev != nil {
		query = "UPDATE jobs SET title=?,category_id=?,expires_at=?,status=? WHERE partner_id=?"
		_, err = jb.db.Query(query, job.Title, job.CategoryID, job.ExpiresAt, job.Status, job.PartnerID)
	} else {
		query = "INSERT INTO jobs (partner_id,title,category_id,expires_at,status) VALUES (?,?,?,?,?)"
		_, err = jb.db.Query(query, job.PartnerID, job.Title, job.CategoryID, job.ExpiresAt, job.Status)
	}
	if err != nil {
		return errors.Wrap(err, "error saving job")
	}

	return nil
}

//Delete erases a job
func (jb *Job) Delete(ID int64) error {
	query := "DELETE FROM jobs WHERE partner_id=?"
	_, err := jb.db.Query(query, ID)
	return err
}

//List retrieves all jobs based on the given pagging and filtering
func (jb *Job) List(limit, page int, status string) ([]*model.Job, error) {
	results := []*model.Job{}
	query := "SELECT partner_id,title,category_id,expires_at,status FROM jobs"
	switch status {
	case model.StatusActive:
		query = fmt.Sprintf("%s WHERE `status`=\"%s\"", query, model.StatusActive)
	case model.StatusDraft:
		query = fmt.Sprintf("%s WHERE `status`=\"%s\"", query, model.StatusDraft)
	}
	skip := limit * page
	query = fmt.Sprintf("%s ORDER BY partner_id LIMIT %d,%d", query, skip, limit)
	rows, err := jb.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		res := model.Job{}
		err = rows.Scan(&res.PartnerID, &res.Title, &res.CategoryID, &res.ExpiresAt, &res.Status)
		if err != nil {
			return nil, errors.Wrap(err, "error retrieving job")
		}
		err = res.Decrypt(jb.cipher)
		if err != nil {
			return nil, errors.Wrap(err, "error retrieving job")
		}
		results = append(results, &res)
	}
	return results, err
}

//Count returns the total amount of jobs with a given status
func (jb *Job) Count(status string) (int, error) {
	var count int
	query := "SELECT COUNT(*) as count FROM jobs"
	switch status {
	case model.StatusActive:
		query = fmt.Sprintf("%s WHERE status=\"%s\"", query, model.StatusActive)
	case model.StatusDraft:
		query = fmt.Sprintf("%s WHERE status=\"%s\"", query, model.StatusDraft)
	}
	rows, err := jb.db.Query(query)

	if err != nil {
		return 0, err
	}
	rows.Next()
	err = rows.Scan(&count)
	return count, err
}

//Percentage calculates the percent and total amount of active jobs per category
func (jb *Job) Percentage(category int64) (*model.CategoryPercentage, error) {
	var count, total float64
	var status string

	rows, err := jb.db.Query("SELECT COUNT(*) as total, status FROM jobs WHERE category_id=? GROUP BY status", category)
	if err != nil {
		return nil, err
	}

	res := &model.CategoryPercentage{}
	res.CategoryID = category

	for rows.Next() {
		err = rows.Scan(&count, &status)
		if err != nil {
			return nil, err
		}
		total += count
		if status == model.StatusActive {
			res.Available = int64(count)
		}
	}
	if total > 0 {
		res.Percentage = (100 * float64(res.Available)) / total
		res.Percentage = float64(int(res.Percentage*100)) / 100
	}
	return res, nil
}

//Activate sets the job status as active
func (jb *Job) Activate(ID int64) error {
	found, err := jb.Get(ID)
	if err != nil {
		return err
	}
	if found.Status == model.StatusActive {
		return nil
	}
	query := "UPDATE jobs SET status=? WHERE partner_id=?"
	_, err = jb.db.Query(query, model.StatusActive, ID)
	return err
}
