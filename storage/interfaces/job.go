//Package interfaces defines all the access methods
//that need to be implemented by a stora backend.
package interfaces

import (
	"errors"

	"github.com/EmpregoLigado/code-challenge/model"
)

//Job is the persistency interface of jobs
type Job interface {
	Get(int64) (*model.Job, error)
	Save(*model.Job) error
	Delete(int64) error
	List(limit, page int, status string) ([]*model.Job, error)
	Count(status string) (int, error)
	Percentage(category int64) (*model.CategoryPercentage, error)
	Activate(int64) error
}

//ErrJobNotFound is returned whenever an action cannot find a job
var ErrJobNotFound = errors.New("job not found")
