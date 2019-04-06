package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestUpdatePost tests that an updated posts
// gets a Now() SyncDate when updated
func TestUpdatePost(t *testing.T) {

	hourago := time.Now().Add(time.Hour * -1)

	// create test server and respond with empty JSON
	// testing the date is set by sync, response from
	// server does not matter
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)

	}))
	defer ts.Close()

	conf.SiteURL = ts.URL

	// create a post in updated array that was
	// previously sync an hourago
	var updatedPosts = []Post{
		Post{
			SyncDate: hourago,
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
