// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestChecksService_GetCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	acceptHeaders := []string{mediaTypeCheckRunsPreview}
	mux.HandleFunc("/repos/o/r/check-runs/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(acceptHeaders, ", "))
		fmt.Fprint(w, `{
			"id": 1,
                        "name":"testCheckRun",
			"status": "completed",
			"conclusion": "neutral",
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": "2018-05-04T01:14:52Z"}`)
	})
	checkRun, _, err := client.Checks.GetCheckRun(context.Background(), "o", "r", 1)
	if err != nil {
		t.Errorf("Checks.GetCheckRun return error: %v", err)
	}
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	completeAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")

	want := &CheckRun{
		ID:          Int64(1),
		Status:      String("completed"),
		Conclusion:  String("neutral"),
		StartedAt:   &startedAt,
		CompletedAt: &completeAt,
		Name:        String("testCheckRun"),
	}
	if !reflect.DeepEqual(checkRun, want) {
		t.Errorf("Checks.GetCheckRun return %+v, want %+v", checkRun, want)
	}
}

func TestChecksService_CreateCheckRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	acceptHeaders := []string{mediaTypeCheckRunsPreview}
	mux.HandleFunc("/repos/o/r/check-runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", strings.Join(acceptHeaders, ", "))
		fmt.Fprint(w, `{
			"id": 1,
                        "name":"testCreateCheckRun",
                        "head_sha":"deadbeef",
			"status": "in_progress",
			"conclusion": null,
			"started_at": "2018-05-04T01:14:52Z",
			"completed_at": null,
                        "output":{"title": "Mighty test report", "summary":"", "text":""}}`)
	})
	startedAt, _ := time.Parse(time.RFC3339, "2018-05-04T01:14:52Z")
	checkRunOpt := &CreateCheckRunOptions{
		HeadBranch: String("master"),
		Name:       String("testCreateCheckRun"),
		HeadSHA:    String("deadbeef"),
		Status:     String("in_progress"),
		StartedAt:  &startedAt,
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String(""),
			Text:    String(""),
		},
	}

	checkRun, _, err := client.Checks.CreateCheckRun(context.Background(), "o", "r", checkRunOpt)
	if err != nil {
		t.Errorf("Checks.CreateCheckRun return error: %v", err)
	}

	want := &CheckRun{
		ID:        Int64(1),
		Status:    String("in_progress"),
		StartedAt: &startedAt,
		HeadSHA:   String("deadbeef"),
		Name:      String("testCreateCheckRun"),
		Output: &CheckRunOutput{
			Title:   String("Mighty test report"),
			Summary: String(""),
			Text:    String(""),
		},
	}
	if !reflect.DeepEqual(checkRun, want) {
		t.Errorf("Checks.CreateCheckRun return %+v, want %+v", checkRun, want)
	}
}
