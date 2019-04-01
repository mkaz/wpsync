package main

import (
	"encoding/json"
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
func uploadPost(filename string) (post Post) {

	page := readParseFile(filename)

	j := getApiFetcher("posts/new")
	j.Params.Add("title", page.Title)
	j.Params.Add("date", page.Date.Format(time.RFC3339))
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)
	j.Params.Add("categories", page.Category)
	j.Params.Add("publicize", "0")
	j.Params.Add("tags", page.Tags)

	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Warn("API error uploading", filename, err)
	}

	if err := json.Unmarshal(resp.Bytes, &post); err != nil {
		log.Warn("Error parsing JSON response", string(resp.Bytes), err)
	}

	return post
}

// create new post
func updatePost(p Post) Post {

	page := readParseFile(p.LocalFile)
	api := fmt.Sprintf("posts/%v", p.Id)
	j := getApiFetcher(api)
	j.Params.Add("title", page.Title)
	j.Params.Add("date", page.Date.Format(time.RFC3339))
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)
	j.Params.Add("categories", page.Category)
	j.Params.Add("publicize", "0")
	j.Params.Add("tags", page.Tags)

	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Warn("API error uploading", p.LocalFile, err)
	}

	if err := json.Unmarshal(resp.Bytes, &p); err != nil {
		log.Warn("Error parsing JSON response", string(resp.Bytes), err)
	}

	return p
}

// upload a single file
func uploadMedia(media Media) Media {

	var ur struct {
		Media []Media
	}

	j := getApiFetcher("media/new")
	j.Files["media[]"] = filepath.Join("media", media.LocalFile)
	resp, err := j.Method("POST").Send()
	if err != nil {
		log.Warn("API error uploading", media.LocalFile, err)
	}

	if err := json.Unmarshal(resp.Bytes, &ur); err != nil {
		log.Warn("Error parsing JSON response", string(resp.Bytes), err)
	}

	if len(ur.Media) > 0 {
		media.Link = ur.Media[0].Link
		media.Id = ur.Media[0].Id
	} else {
		log.Warn("Error: no Link in results", media.LocalFile)
	}
	return media
}
