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
	apiurl := "https://public-api.wordpress.com/rest/v1"
	url := strings.Join([]string{apiurl, "sites", conf.BlogID, endpoint}, "/")

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
		fmt.Println(">>Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &post); err != nil {
		fmt.Println("Error parsing: {} \n\n {}", resp.Bytes, err)
	}

	return post
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
		fmt.Println(">>API Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &ur); err != nil {
		fmt.Println("Error parsing:", err)
		fmt.Println("JSON: %v", string(resp.Bytes))
	}

	if len(ur.Media) > 0 {
		media.Link = ur.Media[0].Link
		media.Id = ur.Media[0].Id
	} else {
		fmt.Println("Error: No Link in results")
		fmt.Println(resp.StatusCode)
	}
	return media
}
