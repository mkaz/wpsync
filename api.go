package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/automattic/go/jaguar"
)

// struct for reading local file in
type Page struct {
	Title, Content, Category, Status, Tags string
	Date                                   time.Time
}

func getApiFetcher(endpoint string) (j jaguar.Jaguar) {
	url := strings.Join([]string{conf.SiteURL, "wp-json", endpoint}, "/")
	j = jaguar.New()
	j.Header.Add("Authorization", "Bearer "+conf.Token)
	j.Url(url)
	return j
}

// create new post
func createPost(filename string) (post Post, err error) {

	page := readParseFile(filename)

	j := getApiFetcher("wp/v2/posts")
	j.Params.Add("title", page.Title)
	j.Params.Add("date", page.Date.Format(time.RFC3339))
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)
	j.Params.Add("publicize", "0")

	resp, err := j.Method("POST").Send()
	if err != nil {
		return post, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("API Error [%v]: %v", resp.StatusCode, string(resp.Bytes))
		return post, errors.New(errMsg)
	}

	err = json.Unmarshal(resp.Bytes, &post)
	return post, err
}

// create new post
func updatePost(p Post) (post Post, err error) {

	page := readParseFile(p.LocalFile)
	api := fmt.Sprintf("wp/v2/posts/%v", p.Id)
	j := getApiFetcher(api)
	j.Params.Add("title", page.Title)
	j.Params.Add("date", page.Date.Format(time.RFC3339))
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)
	j.Params.Add("publicize", "0")

	resp, err := j.Method("POST").Send()
	if err != nil {
		return post, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("API Error [%v]: %v", resp.StatusCode, string(resp.Bytes))
		return post, errors.New(errMsg)
	}

	err = json.Unmarshal(resp.Bytes, &post)
	return post, err
}

// upload a single file
func uploadMedia(media Media) (m Media, err error) {

	j := getApiFetcher("wp/v2/media")
	j.Files["file"] = filepath.Join("media", media.LocalFile)
	resp, err := j.Method("POST").Send()
	if err != nil {
		return m, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("API Error [%v]: %v", resp.StatusCode, string(resp.Bytes))
		return m, errors.New(errMsg)
	}
	err = json.Unmarshal(resp.Bytes, &m)
	if err != nil {
		return m, err
	}

	media.URL = m.URL
	media.Link = m.Link
	media.Id = m.Id
	return media, nil
}
