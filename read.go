package main

import (
	"encoding/json"
)

func get_latest() []Post {
	posts := parseFetchPosts()
	return posts
}

// fetch single post
func get_single_post(post_id string) {
	post := parseFetchPost(post_id)
	log.Debug("\n# " + post.Title + "\n")
	log.Debug(post.Content)
	log.Debug(post.URL + "\n")
}

// fetch and parse list of posts
func parseFetchPosts() []Post {
	j := getApiFetcher("posts/")
	resp, err := j.Send()
	if err != nil {
		log.Fatal(">>Error: ", err)
	}

	var h []Post
	if err := json.Unmarshal(resp.Bytes, &h); err != nil {
		log.Fatal("Error parsing:", err)
	}

	return h
}

// parse single post
func parseFetchPost(post_id string) (p Post) {
	j := getApiFetcher("posts/" + post_id)
	resp, err := j.Send()
	if err != nil {
		log.Fatal(">>Error: ", err)
	}

	if err := json.Unmarshal(resp.Bytes, &p); err != nil {
		log.Fatal("Error parsing:", err)
	}
	return p
}
