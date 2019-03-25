package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// getLocalPosts reads posts from local directory
func getLocalPosts() (posts []Post) {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal("Error reading directory", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".md") {
			post := Post{}
			fmt.Println("Filename: " + file.Name())
			post.LocalFile = file.Name()
			posts = append(posts, post)
		}
	}
	return posts
}

// getRemotePosts reads posts from json file
func getRemotePosts() (posts []Post) {
	// TODO: read json file

	posts = append(posts, Post{LocalFile: "sample-post.md"})
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

func uploadPosts(posts []Post) []Post {
	for _, p := range posts {
		uploadPost(p.LocalFile)
	}
	return posts
}

// getLocalMedia reads media from local directory
func getLocalMedia() (media []Media) {
	files, err := ioutil.ReadDir("./media")
	if err != nil {
		log.Fatal("Error reading directory", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".jpg") {
			m := Media{}
			fmt.Println("Media file: " + file.Name())
			m.LocalFile = file.Name()
			media = append(media, m)
		}
	}
	return media
}

// getRemoteMedia reads media from json file
func getRemoteMedia() (media []Media) {
	return media
}

func compareMedia(local, remote []Media) (media []Media) {
	for _, m := range local {
		exists := false
		for _, r := range remote {
			if m.LocalFile == r.LocalFile {
				exists = true
			}
		}
		if !exists {
			media = append(media, m)
		}
	}
	return media
}

func uploadMedia(media []Media) []Media {
	return media
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
		log.Fatalln(">>Error: can't read file:", filename)
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
