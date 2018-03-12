package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/EmpregoLigado/code-challenge/model"
	"github.com/EmpregoLigado/code-challenge/storage/interfaces"
)

type Backend struct {
	Name    string
	Backend interfaces.Job
}

var jobs = []*model.Job{
	&model.Job{
		PartnerID:  1,
		CategoryID: 1,
		Title:      "First Job",
		Status:     model.StatusActive,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	},
	&model.Job{
		PartnerID:  2,
		CategoryID: 1,
		Title:      "Second Job",
		Status:     model.StatusAny,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	},
	&model.Job{
		PartnerID:  3,
		CategoryID: 2,
		Title:      "Third Job",
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	},
	&model.Job{
		PartnerID:  4,
		CategoryID: 3,
		Title:      "Fourth Job",
		Status:     model.StatusDraft,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
	},
}

var xjob = &model.Job{
	PartnerID:  5,
	CategoryID: 1,
	Title:      "Extra Job",
	Status:     model.StatusActive,
	ExpiresAt:  time.Now().Add(1 * time.Hour),
}

//JobBackendTest utility function to help test the job backends
func JobBackendTest(t *testing.T, backends []Backend) {
	for _, backend := range backends {
		sto := backend.Backend
		t.Run(fmt.Sprintf("%s job save", backend.Name), func(t *testing.T) {
			for _, job := range jobs {
				err := sto.Save(job)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
		t.Run(fmt.Sprintf("%s job count", backend.Name), func(t *testing.T) {
			count, err := sto.Count(model.StatusDraft)
			if err != nil {
				t.Fatal(err)
			}
			if count != len(jobs) {
				t.Fatalf("should have returned a count of %d", len(jobs))
			}
		})
		t.Run(fmt.Sprintf("%s job status", backend.Name), func(t *testing.T) {
			for _, job := range jobs {
				jb, err := sto.Get(job.PartnerID)
				if err != nil {
					t.Fatal(err)
				}
				if jb.Status != model.StatusDraft {
					t.Fatal("every new job must have a draft status")
				}
			}
		})
		t.Run(fmt.Sprintf("%s job activate", backend.Name), func(t *testing.T) {
			for _, job := range jobs {
				err := sto.Activate(job.PartnerID)
				if err != nil {
					t.Fatal(err)
				}
				jb, err := sto.Get(job.PartnerID)
				if err != nil {
					t.Fatal(err)
				}
				if jb.Status != model.StatusActive {
					t.Fatal("activated job must have an active status")
				}
			}
		})
		t.Run(fmt.Sprintf("%s job percentage", backend.Name), func(t *testing.T) {
			percentage, err := sto.Percentage(1)
			if err != nil {
				t.Fatal(err)
			}
			if percentage.Available != 2 {
				t.Fatalf("should have 2 elements in percentage total, returned %d", percentage.Available)
			}
			if percentage.Percentage != 100.0 {
				t.Fatalf("should have a percentage of 100%%, got %f", percentage.Percentage)
			}
			err = sto.Save(xjob)
			if err != nil {
				t.Fatal(err)
			}

			percentage, err = sto.Percentage(1)
			if err != nil {
				t.Fatal(err)
			}
			if percentage.Available != 2 {
				t.Fatalf("should have 2 elements in percentage total, returned %d", percentage.Available)
			}
			if percentage.Percentage != 66.66 {
				t.Fatalf("should have a percentage of 66.66%%, got %f", percentage.Percentage)
			}

			percentage, err = sto.Percentage(999)
			if err != nil {
				t.Fatal(err)
			}
			if percentage.Available != 0 {
				t.Fatalf("should have 0 elements in percentage total, returned %d", percentage.Available)
			}
			if percentage.Percentage != 0 {
				t.Fatalf("should have a percentage of 0%%, got %f", percentage.Percentage)
			}

		})
		t.Run(fmt.Sprintf("%s job list", backend.Name), func(t *testing.T) {
			//Unlimited
			jbs, err := sto.List(len(jobs)+1, 0, model.StatusActive)
			if err != nil {
				t.Fatal(err)
			}
			if len(jbs) != len(jobs) {
				t.Fatalf("should have returned %d jobs, returned %d", len(jobs), len(jbs))
			}
			//First element
			jbs, err = sto.List(1, 0, model.StatusActive)
			if err != nil {
				t.Fatal(err)
			}
			if len(jbs) != 1 {
				t.Fatalf("should have returned %d jobs", 1)
			}
			//Last element
			jbs, err = sto.List(1, len(jobs)-1, model.StatusActive)
			if err != nil {
				t.Fatal(err)
			}
			if len(jbs) != 1 {
				t.Fatalf("should have returned %d jobs", 1)
			}
			//First two elements
			jbs, err = sto.List(2, 0, model.StatusActive)
			if err != nil {
				t.Fatal(err)
			}
			if len(jbs) != 2 {
				t.Fatalf("should have returned %d jobs", 2)
			}
			//Last two elements
			jbs, err = sto.List(2, (len(jobs)/2)-1, model.StatusActive)
			if err != nil {
				t.Fatal(err)
			}
			if len(jbs) != 2 {
				t.Fatalf("should have returned %d jobs", 2)
			}
		})
		t.Run(fmt.Sprintf("%s job delete", backend.Name), func(t *testing.T) {
			err := sto.Delete(xjob.PartnerID)
			if err != nil {
				t.Fatal(err)
			}
			_, err = sto.Get(xjob.PartnerID)
			if err == nil {
				t.Fatal("should have failed with not found error")
			}

			for _, job := range jobs {
				err := sto.Delete(job.PartnerID)
				if err != nil {
					t.Fatal(err)
				}
				_, err = sto.Get(job.PartnerID)
				if err == nil {
					t.Fatal("should have failed with not found error")
				}
			}
		})
	}
}
