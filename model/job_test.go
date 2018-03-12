package model

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestOutputJobToJson(t *testing.T) {
	dts := "19/6/1979"
	dt, _ := time.Parse(DateFormat, dts)
	job := OutputJob{
		PartnerID:  1009,
		CategoryID: 9001,
		Title:      "Fake Job One",
		ExpiresAt:  dt,
	}
	jsn, err := json.Marshal(&job)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(jsn), dts) {
		t.Fatalf("could not find date string %s, json: %s", dts, string(jsn))
	}
}
