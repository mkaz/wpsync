// wpsync - command-line tool to sync wordpress
// https://github.com/mkaz/wpsync
//
// TODO: add watch
// TODO: add confirmation flag
// TODO: add download flag

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type Config struct {
	SiteURL string `json:"site-url"`
	Token   string `json:"token"   `
	Dryrun  bool
}

type Post struct {
	Id        int    `json:"id"`
	Title     string `json:"title.raw"`
	Date      WPTime `json:"date"`
	URL       string `json:"URL"`
	Content   string `json:"content.raw"`
	LocalFile string
}

type Media struct {
	Id        int    `json:"id"`
	URL       string `json:"source_url"`
	Link      string `json:"link"`
	LocalFile string
}

var conf Config
var log Logger
var setup bool

// read config and parse args
func init() {

	flag.BoolVar(&log.Verbose, "verbose", false, "Chatty")
	flag.BoolVar(&conf.Dryrun, "dryrun", false, "No uploads")
	flag.BoolVar(&setup, "init", false, "Setup and Test")
	flag.Parse()

	file, err := ioutil.ReadFile("wpsync.json")
	if err != nil {
		log.Debug("wpsync.json file not found, running setup", err)
		setup = true
	} else {
		if err := json.Unmarshal(file, &conf); err != nil {
			log.Fatal("Error parsing wpsync.json", err)
		}
		log.Debug("Config loaded", conf)
	}

	if setup {
		runSetup()
	}

	// test setup
	if !testSetup() {
		// setup not working
		// check if runSetup() ran with setup
		// if not run it now otherwise bail
		log.Fatal("Setup not confirmed.", conf)
	}

}

// route command and args
func main() {

	// read local files for data
	localPosts := getLocalPosts()
	log.Debug("Found local posts: ", localPosts)

	remotePosts := getRemotePosts()
	log.Debug("Existing posts: ", remotePosts)

	newPosts, updatedPosts := comparePosts(localPosts, remotePosts)
	log.Debug("New posts to upload: ", newPosts)
	log.Debug("Existing post to update: ", updatedPosts)

	if !conf.Dryrun {
		newPosts = uploadPosts(newPosts)
		updatedPosts = updatePosts(updatedPosts)
		writeRemotePosts(newPosts)
		for _, p := range newPosts {
			log.Info("New Post: ", p.LocalFile, p.URL)
		}
	}

	// media
	localMedia := getLocalMedia()
	for _, m := range localMedia {
		log.Debug("Found local media: ", m.LocalFile)
	}

	remoteMedia := getRemoteMedia()
	for _, m := range remoteMedia {
		log.Debug("Existing media: ", m.LocalFile)
	}

	newMedia := compareMedia(localMedia, remoteMedia)
	for _, m := range newMedia {
		log.Debug("New media to upload: ", m.LocalFile)
	}

	if !conf.Dryrun {
		uploadMediaItems(newMedia)
		writeRemoteMedia(newMedia)
		for _, m := range newMedia {
			log.Info(m.LocalFile, m.Link)
		}
	}

}
