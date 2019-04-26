package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/automattic/go/jaguar"
)

func getApiFetcher(endpoint string) (j jaguar.Jaguar) {
	url := strings.Join([]string{conf.SiteURL, "wp-json", endpoint}, "/")
	j = jaguar.New()
	j.Header.Add("Authorization", "Bearer "+conf.Token)
	j.Url(url)
	return j
}

// create new post
func createPost(post Post) (Post, error) {
	j := getApiFetcher("wp/v2/posts")
	j.Params.Add("title", post.Title)
	j.Params.Add("date", post.Date)
	j.Params.Add("content", post.Content)
	j.Params.Add("status", post.Status)
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

// update existing post
func updatePost(post Post) (Post, error) {

	api := fmt.Sprintf("wp/v2/posts/%v", post.Id)
	j := getApiFetcher(api)
	j.Params.Add("title", post.Title)
	j.Params.Add("date", post.Date)
	j.Params.Add("content", post.Content)
	j.Params.Add("status", post.Status)
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

// create new page
func createPage(page Page) (Page, error) {
	j := getApiFetcher("wp/v2/pages")
	j.Params.Add("title", page.Title)
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)

	if page.Template != "" {
		j.Params.Add("template", page.Template)
	}

	if page.ParentId != 0 {
		j.Params.Add("parent", strconv.Itoa(page.ParentId))
	}

	if page.Order != "" {
		j.Params.Add("menu_order", page.Order)
	}

	resp, err := j.Method("POST").Send()
	log.Debug("Making request", string(resp.Bytes))
	if err != nil {
		return page, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("API Error [%v]: %v", resp.StatusCode, string(resp.Bytes))
		return page, errors.New(errMsg)
	}

	err = json.Unmarshal(resp.Bytes, &page)
	return page, err
}

// update existing page
func updatePage(page Page) (Page, error) {

	api := fmt.Sprintf("wp/v2/pages/%v", page.Id)
	j := getApiFetcher(api)
	j.Params.Add("title", page.Title)
	j.Params.Add("content", page.Content)
	j.Params.Add("status", page.Status)

	if page.Template != "" {
		j.Params.Add("template", page.Template)
	}

	if page.ParentId != 0 {
		j.Params.Add("parent", strconv.Itoa(page.ParentId))
	}

	if page.Order != "" {
		j.Params.Add("menu_order", page.Order)
	}

	resp, err := j.Method("POST").Send()
	if err != nil {
		return page, err
	}

	if resp.StatusCode > 299 {
		errMsg := fmt.Sprintf("API Error [%v]: %v", resp.StatusCode, string(resp.Bytes))
		return page, errors.New(errMsg)
	}

	err = json.Unmarshal(resp.Bytes, &page)
	return page, err
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
