package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestPublishStatus gets set properly (publish/draft)
func TestPublishStatus(t *testing.T) {

	// Server to echo status back, or publish (wp default) if not set
	statusHandler := func(w http.ResponseWriter, r *http.Request) {
		status := r.FormValue("status")
		if status == "" {
			status = "publish"
		}

		post := Post{
			Status: status,
		}

		js, _ := json.Marshal(post)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	ts := httptest.NewServer(http.HandlerFunc(statusHandler))
	defer ts.Close()

	conf.SiteURL = ts.URL

	var newPosts = []Post{
		Post{
			LocalFile: "draft.md",
			Status:    "draft",
		},
		Post{
			LocalFile: "publish.md",
			Status:    "publish",
		},
	}

	createdPosts := createPosts(newPosts)
	// check createPosts status
	if len(createdPosts) == 0 {
		t.Error("No created posts")
	} else {
		if createdPosts[0].Status != "draft" {
			t.Error("New draft post not draft status")
		}

		if createdPosts[1].Status != "publish" {
			t.Error("New publish post not publish status")
		}
	}
}

// TestUpdatePost tests an updated posts gets an updated SyncDate
func TestUpdatePost(t *testing.T) {

	hourago := time.Now().Add(time.Hour * -1)

	emptyHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	}

	// create test server and respond with empty JSON
	// testing the date is set by sync, response from
	// server does not matter
	ts := httptest.NewServer(http.HandlerFunc(emptyHandler))
	defer ts.Close()

	conf.SiteURL = ts.URL

	// create a post in updated array that was
	// previously sync an hourago
	var updatedPosts = []Post{
		Post{
			LocalFile: "hourago.md",
			SyncDate:  hourago,
		},
	}

	// Do the thing
	updatedPosts = updatePosts(updatedPosts)

	// confirm post sync date is updated
	if len(updatedPosts) == 0 {
		t.Error("No updated posts")
	} else {
		if updatedPosts[0].SyncDate.Before(hourago) {
			t.Error("Updated post sync date not updated")
		}
	}
}
