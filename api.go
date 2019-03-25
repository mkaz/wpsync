package main

import (
	"encoding/json"
	"fmt"
	"log"
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
func uploadPost(filename string) {

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
		log.Fatalln(">>Error: ", err)
	}

	newurl := parseNewPostResponse(resp.Bytes)

	fmt.Println("New Post:", newurl)
}

// extract URL from json response data of new post
func parseNewPostResponse(data []byte) string {

	var rs struct{ Url string }

	if err := json.Unmarshal(data, &rs); err != nil {
		log.Fatalf("Error parsing: {} \n\n {}", data, err)
	}

	return rs.Url
}

// upload a single file
func upload_media(filename string) {
	var ur struct {
		Media []struct {
			Link  string
			Title string
		}
	}

	j := getApiFetcher("media/new")
	j.Files["media[]"] = filename
	resp, err := j.Method("POST").Send()

	if err != nil {
		log.Fatalln(">>Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &ur); err != nil {
		log.Fatal("Error parsing:", err)
	}

	if len(ur.Media) > 0 {
		fmt.Println(ur.Media[0].Link)
	} else {
		fmt.Println("Error: No link in results")
		fmt.Println(resp.StatusCode)
	}

}
