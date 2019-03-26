package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// getLocalPosts reads posts from local directory
func getLocalPosts() (posts []Post) {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal("Error reading posts directory", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".md") {
			post := Post{}
			post.LocalFile = file.Name()
			posts = append(posts, post)
		}
	}
	return posts
}

// getRemotePosts reads posts from json file
func getRemotePosts() (posts []Post) {
	// check if file exists, return empty
	// likely scenario would be first run
	if _, err := os.Stat("posts.json"); os.IsNotExist(err) {
		fmt.Println("INFO: posts.json does not exist, first run?")
		return posts
	}

	file, err := ioutil.ReadFile("posts.json")
	if err != nil {
		fmt.Println("Error reading posts.json, permissions?", err)
	} else {
		if err := json.Unmarshal(file, &posts); err != nil {
			fmt.Println("Error parsing JSON from posts.json", err)
		}
	}
	return posts
}

// comparePosts returns local posts that do not exist in remote
func comparePosts(local, remote []Post) (posts []Post) {
	for _, p := range local {
		exists := false
		for _, r := range remote {
			if p.LocalFile == r.LocalFile {
				exists = true
				fmt.Println("Skipping ", p.LocalFile)
			}
		}
		if !exists {
			posts = append(posts, p)
		}
	}
	return posts
}

// uploadPosts loops through posts and uploads
// posts are returned with Id/Url set
func uploadPosts(posts []Post) []Post {
	for i, p := range posts {
		p = uploadPost(p.LocalFile)
		posts[i].Id = p.Id
		posts[i].URL = p.URL
	}
	return posts
}

// writeRemotePosts
func writeRemotePosts(posts []Post) {
	if len(posts) == 0 {
		fmt.Println("No new posts to write.")
		return
	}
	// append new post json
	existingPosts := getRemotePosts()
	existingPosts = append(existingPosts, posts...)

	// write file
	json, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("JSON Encoding Error", err)
	} else {
		err = ioutil.WriteFile("posts.json", json, 0644)
		if err != nil {
			fmt.Println("Error writing posts.json", err)
		} else {
			fmt.Println("posts.json written")
		}
	}
}

// readParseFile reads a markdown file and returns a page struct
func readParseFile(filename string) (page Page) {

	// setup default page struct
	page = Page{
		Title:    "",
		Content:  "",
		Category: "",
		Date:     time.Now(),
		Tags:     "",
		Status:   "publish",
	}

	var data, err = ioutil.ReadFile(filepath.Join("posts", filename))
	if err != nil {
		log.Fatal(">>Error: can't read file:", filename)
	}

	// parse front matter from --- to ---
	var lines = strings.Split(string(data), "\n")
	var found = 0
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if found == 1 {
			// parse line for param
			colonIndex := strings.Index(line, ":")
			if colonIndex > 0 {
				key := strings.TrimSpace(line[:colonIndex])
				value := strings.TrimSpace(line[colonIndex+1:])
				value = strings.Trim(value, "\"") //remove quotes
				switch key {
				case "title":
					page.Title = value
				case "date":
					page.Date, _ = time.Parse("2006-01-02", value)
				case "category":
					page.Category = value
				case "tags":
					page.Tags = value
				case "status":
					page.Status = value
				}
			}
		} else if found >= 2 {
			// params over
			lines = lines[i:]
			break
		}

		if line == "---" {
			found += 1
		}
	}

	// slurp rest of content
	page.Content = strings.Join(lines, "\n")
	return page
}
