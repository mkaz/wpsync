package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"gopkg.in/russross/blackfriday.v2"
)

// getLocalPosts reads posts from local directory
func getLocalPosts() (posts []Post) {
	files, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Info("Error reading posts directory: %v", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".md") {
			post := Post{}
			post.LocalFile = file.Name()
			post.ModDate = file.ModTime()
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
		if !setup { // dont alert about missing file when known init
			log.Debug("posts.json does not exist")
		}
		return posts
	}

	file, err := ioutil.ReadFile("posts.json")
	if err != nil {
		log.Warn("Error reading posts.json, permissions?", err)
	} else {
		if err := json.Unmarshal(file, &posts); err != nil {
			log.Warn("Error parsing JSON from posts.json", err)
		}
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
				if lp.ModDate.After(rp.SyncDate) {
					log.Debug("Local File: ", lp.LocalFile)
					log.Debug("   Local ModDate  : ", lp.ModDate.Unix())
					log.Debug("   Remote SyncDate: ", rp.SyncDate.Unix())
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

// createPosts loops through posts and uploads
// posts are returned with Id/Url set
func createPosts(newPosts []Post) (createdPosts []Post) {
	for _, p := range newPosts {
		if confirmPrompt(fmt.Sprintf("New post %s, Continue (y/N)? ", p.LocalFile)) {
			rp, err := createPost(p)
			if err == nil {
				rp.LocalFile = p.LocalFile // do I need to merge all data
				rp.SyncDate = time.Now()
				log.Info(fmt.Sprintf("New post: %s %s", p.LocalFile, rp.URL))
				createdPosts = append(createdPosts, rp)
			}
		}
	}
	return createdPosts
}

func loadPostsFromFiles(posts []Post) (loadedPosts []Post) {
	for _, p := range posts {
		lp := loadPostFromFile(p)
		loadedPosts = append(loadedPosts, lp)
	}
	return loadedPosts
}

func loadPostFromFile(p Post) Post {
	post := readParseFile(p.LocalFile)
	mergo.Merge(&post, p)
	return post
}

// updatePosts loops through posts and updates
// posts are returned with new Date set
func updatePosts(posts []Post) (updatedPosts []Post) {
	for _, p := range posts {
		if confirmPrompt(fmt.Sprintf("Update post %s, Continue (y/N)? ", p.LocalFile)) {
			rp, err := updatePost(p)
			if err == nil {
				rp.SyncDate = time.Now()
				log.Info(fmt.Sprintf("Updated post: %s %s", p.LocalFile, rp.URL))
				log.Debug("Updated SyncDate to:", rp.SyncDate.Unix())
				updatedPosts = append(updatedPosts, rp)
			} else {
				log.Warn("Error updating post", err)
			}
		}
	}
	return updatedPosts
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
				existingPosts[i].SyncDate = up.SyncDate
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

// readParseFile reads a markdown file and returns a Post struct
func readParseFile(filename string) (post Post) {

	// setup default data
	post = Post{
		Title:    "",
		Content:  "",
		Category: "",
		Date:     time.Now().Format(time.RFC3339),
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
					post.Title = value
				case "date":
					d, err := time.Parse("2006-01-02", value)
					if err == nil {
						post.Date = d.Format(time.RFC3339)
					}
				case "category":
					post.Category = value
				case "tags":
					post.Tags = value
				case "status":
					post.Status = value
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
	post.Content = string(blackfriday.Run([]byte(content)))

	return post
}
