package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"gopkg.in/russross/blackfriday.v2"
)

// getLocalPages reads  from local directory
func getLocalPages() (pages []Page) {
	files, err := ioutil.ReadDir("./pages")
	if err != nil {
		log.Info("Error reading pages directory: %v", err)
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".md") {
			log.Debug("Pages file name:", file.Name())
			page := Page{}
			page.LocalFile = file.Name()
			page.ModDate = file.ModTime()
			pages = append(pages, page)
		}
	}
	return pages
}

// getRemotePages reads pages from json file
func getRemotePages() (pages []Page) {
	// check if file exists, return empty
	// likely scenario would be first run
	if _, err := os.Stat("pages.json"); os.IsNotExist(err) {
		if !setup { // dont alert about missing file when known init
			log.Debug("pages.json does not exist")
		}
		return pages
	}

	file, err := ioutil.ReadFile("pages.json")
	if err != nil {
		log.Warn("Error reading pages.json, permissions?", err)
	} else {
		if err := json.Unmarshal(file, &pages); err != nil {
			log.Warn("Error parsing JSON from pages.json", err)
		}
	}
	return pages
}

// comparePages returns local pages that do not exist in remote
func comparePages(local, remote []Page) (newPages, updatePages []Page) {
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
					updatePages = append(updatePages, lp)
				} else {
					log.Debug("Skipping ", lp.LocalFile)
				}
			}
		}
		if !exists {
			newPages = append(newPages, lp)
		}
	}
	return newPages, updatePages
}

// createPages loops through pages and uploads
// pages are returned with Id/Url set
func createPages(newPages []Page) (createdPages []Page) {
	for _, p := range newPages {
		if confirmPrompt(fmt.Sprintf("New page %s, Continue (y/N)? ", p.LocalFile)) {
			rp, err := createPage(p)
			if err == nil {
				rp.LocalFile = p.LocalFile // do I need to merge all data
				rp.SyncDate = time.Now()
				log.Info(fmt.Sprintf("New page: %s %s", p.LocalFile, rp.URL))
				createdPages = append(createdPages, rp)
			} else {
				log.Debug("Error creating page", err)
			}
		}
	}
	return createdPages
}

func loadPagesFromFiles(pages []Page) (loadedPages []Page) {
	for _, p := range pages {
		lp := loadPageFromFile(p)
		loadedPages = append(loadedPages, lp)
	}
	return loadedPages
}

func loadPageFromFile(p Page) Page {
	page := readParsePageFile(p.LocalFile)
	mergo.Merge(&page, p)
	return page
}

// updatePages loops through pages and updates
// pages are returned with new Date set
func updatePages(pages []Page) (updatedPages []Page) {
	for _, p := range pages {
		if confirmPrompt(fmt.Sprintf("Update page %s, Continue (y/N)? ", p.LocalFile)) {
			rp, err := updatePage(p)
			if err == nil {
				rp.SyncDate = time.Now()
				log.Info(fmt.Sprintf("Updated page: %s %s", p.LocalFile, rp.URL))
				log.Debug("Updated SyncDate to:", rp.SyncDate.Unix())
				updatedPages = append(updatedPages, rp)
			} else {
				log.Warn("Error updating page", err)
			}
		}
	}
	return updatedPages
}

// writeRemotePages
func writeRemotePages(newPages, updatedPages []Page) {
	if len(newPages) == 0 && len(updatedPages) == 0 {
		log.Info("No pages to write.")
		return
	}
	// append new pages.json
	existingPages := getRemotePages()
	// Merge existingPages and updatedPages
	// need to update the date
	for i, ep := range existingPages {
		for _, up := range updatedPages {
			if ep.LocalFile == up.LocalFile {
				existingPages[i].SyncDate = up.SyncDate
			}
		}
	}
	existingPages = append(existingPages, newPages...)

	// write file
	json, err := json.Marshal(existingPages)
	if err != nil {
		log.Warn("JSON Encoding Error", err)
	} else {
		err = ioutil.WriteFile("pages.json", json, 0644)
		if err != nil {
			log.Warn("Error writing pages.json", err)
		} else {
			log.Debug("pages.json written")
		}
	}
}

// readParseFile reads a markdown file and returns a Page struct
func readParsePageFile(filename string) (page Page) {

	// setup default data
	page = Page{
		Title:    "",
		Content:  "",
		Template: "",
		ParentId: 0,
		Status:   "publish",
	}

	var data, err = ioutil.ReadFile(filepath.Join("pages", filename))
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
				case "template":
					page.Template = value
				case "parent":
					page.ParentId, _ = strconv.Atoi(value)
				case "status":
					page.Status = value
				case "order":
					page.Order = value
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
