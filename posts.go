package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/russross/blackfriday.v2"
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
			post.Date = WPTime{file.ModTime()}
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
		log.Info("posts.json does not exist, first run?")
		return posts
	}

	file, err := ioutil.ReadFile("posts.json")
	if err != nil {
		log.Warn("Error reading posts.json, permissions?", err)
	} else {
		if err := json.Unmarshal(file, &posts); err != nil {
			log.Warn("Error parsing JSON from posts.json", err)
		}
		log.Debug("Posts unmarshal", posts)
	}
	return posts
}

// comparePosts returns local posts that do not exist in remote
func comparePosts(local, remote []Post) (newPosts, updatePosts []Post) {
	for _, lp := range local {
		exists := false
		for _, rp := range remote {
			if lp.LocalFile == rp.LocalFile {
				exists = true
				lp.Id = rp.Id // set Id from remote
				if lp.Date.After(rp.Date.Time) {
					log.Debug("Local Date : ", lp.Date.Unix())
					log.Debug("Remote Date: ", rp.Date.Unix())
					updatePosts = append(updatePosts, lp)
				} else {
					log.Debug("Skipping ", lp.LocalFile)
				}
			}
		}
		if !exists {
			newPosts = append(newPosts, lp)
		}
	}
	return newPosts, updatePosts
}

// uploadPosts loops through posts and uploads
// posts are returned with Id/Url set
func uploadPosts(posts []Post) []Post {
	for i, p := range posts {
		rp := uploadPost(p.LocalFile)
		posts[i].Id = rp.Id
		posts[i].URL = rp.URL

		// Only update if date from remote post is after
		// this makes sure when updating a post the new
		// updated date is used, not the remote post's date
		if posts[i].Date.Before(rp.Date.Time) {
			posts[i].Date = rp.Date
		}
	}
	return posts
}

// udatePosts loops through posts and updates
// posts are returned with new Date set
func updatePosts(posts []Post) []Post {
	for i, p := range posts {
		log.Debug("Updating post", p.Id, p.LocalFile)
		rp := updatePost(p)
		// Only update if date from remote post is after
		// this makes sure when updating a post the new
		// updated date is used, not the remote post's date
		if posts[i].Date.Before(rp.Date.Time) {
			posts[i].Date = rp.Date
		}
	}
	return posts
}

// writeRemotePosts
func writeRemotePosts(newPosts, updatedPosts []Post) {
	if len(newPosts) == 0 && len(updatedPosts) == 0 {
		log.Info("No posts to write.")
		return
	}
	// append new post json
	existingPosts := getRemotePosts()
	// Merge existingPosts and updatedPosts
	// need to update the date
	for i, ep := range existingPosts {
		for _, up := range updatedPosts {
			if ep.LocalFile == up.LocalFile {
				existingPosts[i].Date = up.Date
			}
		}
	}
	existingPosts = append(existingPosts, newPosts...)

	// write file
	json, err := json.Marshal(existingPosts)
	if err != nil {
		log.Warn("JSON Encoding Error", err)
	} else {
		err = ioutil.WriteFile("posts.json", json, 0644)
		if err != nil {
			log.Warn("Error writing posts.json", err)
		} else {
			log.Debug("posts.json written")
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
		log.Warn(">>Error: can't read file:", filename)
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
	content := strings.Join(lines, "\n")
	page.Content = string(blackfriday.Run([]byte(content)))

	return page
}
