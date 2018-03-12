package memory

import (
	"errors"
	"fmt"
	"sort"

	"github.com/EmpregoLigado/code-challenge/crypt"
	"github.com/EmpregoLigado/code-challenge/model"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
)

//Job base type to access in memory jobs
type Job struct {
	cipher crypt.Cipher
	jobs   map[int64]*model.Job
}

//NewJob creates a new instance of a memory job storage
func NewJob(cipher crypt.Cipher) *Job {
	job := &Job{}
	job.jobs = make(map[int64]*model.Job)
	job.cipher = cipher
	return job
}

//Get retrieves one job by partnerID
func (jb *Job) Get(ID int64) (*model.Job, error) {
	r, ok := jb.jobs[ID]
	if ok {
		err := r.Decrypt(jb.cipher)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, interfaces.ErrJobNotFound
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
	job.Status = model.StatusDraft

	err := job.Encrypt(jb.cipher)
	if err != nil {
		return err
	}
	jb.jobs[job.PartnerID] = job
	return nil
}

//Delete erases a job
func (jb *Job) Delete(ID int64) error {
	_, ok := jb.jobs[ID]
	if !ok {
		return fmt.Errorf("job %d does not exist", ID)
	}
	delete(jb.jobs, ID)
	return nil
}

//List retrieves all jobs based on the given pagging and filtering
func (jb *Job) List(limit, page int, status string) ([]*model.Job, error) {
	var res []*model.Job
	var i, j int
	offset := limit * page

	//Ensures that all the results are sorted by partner_id.
	//Since we dont have any sorting request option, this avoids the default
	//go behavior of randomizing map item order.
	var keys []int64
	for k := range jb.jobs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, key := range keys {
		val := jb.jobs[key]
		if status != model.StatusAny && val.Status != status {
			continue
		}
		j++
		if j < offset {
			continue
		}
		err := val.Decrypt(jb.cipher)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
		i++
		if i == limit {
			break
		}
	}
	return res, nil
}

//Count returns the total amount of jobs with a given status
func (jb *Job) Count(status string) (int, error) {
	var count int
	for _, val := range jb.jobs {
		if val.Status != status {
			continue
		}
		count++
	}
	return count, nil
}

//Percentage calculates the percent and total amount of active jobs per category
func (jb *Job) Percentage(category int64) (*model.CategoryPercentage, error) {
	res := &model.CategoryPercentage{}
	res.CategoryID = category
	var count float64
	for _, val := range jb.jobs {
		if val.CategoryID != category {
			continue
		}
		count++
		if val.Status == model.StatusActive {
			res.Available++
		}
	}
	if count > 0 {
		res.Percentage = (100 * float64(res.Available)) / count
		res.Percentage = float64(int(res.Percentage*100)) / 100
	}
	return res, nil
}

//Activate sets the job status as active
func (jb *Job) Activate(ID int64) error {
	job, err := jb.Get(ID)
	if err != nil {
		return err
	}
	job.Status = model.StatusActive
	err = job.Encrypt(jb.cipher)
	if err != nil {
		return err
	}
	jb.jobs[job.PartnerID] = job
	return nil
}
