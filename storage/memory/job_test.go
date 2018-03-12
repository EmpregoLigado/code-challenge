package memory

import (
	"testing"

	"github.com/EmpregoLigado/code-challenge/crypt/null"
	"github.com/EmpregoLigado/code-challenge/storage/test"
)

func TestJobBackend(t *testing.T) {
	backends := []test.Backend{}
	cipher := &null.Cipher{}
	backends = append(backends, test.Backend{Name: "memory", Backend: NewJob(cipher)})
	test.JobBackendTest(t, backends)
}
