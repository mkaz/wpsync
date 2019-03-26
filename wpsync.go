// wpsync - command-line tool to sync wordpress
// https://github.com/mkaz/wpsync
//
// TODO: add watch
// TODO: add confirmation flag
// TODO: add download flag

package main

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BlogID string `toml:"blog_id"`
	Token  string `toml:"token"`
	Dryrun bool
}

type Post struct {
	Id        int       `json:"ID"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	URL       string    `json:"URL"`
	Content   string    `json:"content"`
	LocalFile string
}

type Media struct {
	Id        string    `json:"id"`
	Date      time.Time `json:"date"`
	Link      string    `json:"link"`
	LocalFile string
}

var conf Config
var log Logger

// read config and parse args
func init() {

	configFilename := "wpsync.conf"

	// confirm file exists
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		log.Fatal(">> Config file wpsync.conf does not exists", err)
	}

	// parse file
	if _, err := toml.DecodeFile(configFilename, &conf); err != nil {
		log.Fatal(">> Error decoding wpsync.conf config file", err)
	}

	// confirm params and config all set
	if conf.Token == "" {
		log.Fatal("No auth token configured in wpsync.conf")
	}

	if conf.BlogID == "" {
		log.Fatal("No blog id configured in wpsync.conf")
	}

	flag.BoolVar(&log.Verbose, "verbose", false, "Chatty")
	flag.BoolVar(&conf.Dryrun, "dryrun", false, "No uploads")
	flag.Parse()
}

// route command and args
func main() {

	// read local files for data
	localPosts := getLocalPosts()
	for _, p := range localPosts {
		log.Debug("Found local posts: ", p.LocalFile)
	}

	remotePosts := getRemotePosts()
	for _, p := range remotePosts {
		log.Debug("Existing posts: ", p.LocalFile)
	}

	newPosts := comparePosts(localPosts, remotePosts)
	for _, p := range newPosts {
		log.Debug("New posts to upload: ", p.LocalFile)
	}

	if !conf.Dryrun {
		uploadPosts(newPosts)
		writeRemotePosts(newPosts)
		for _, p := range newPosts {
			log.Info("New Post: ", p.LocalFile, p.URL)
		}
	}

	// media
	localMedia := getLocalMedia()
	remoteMedia := getRemoteMedia()
	newMedia := compareMedia(localMedia, remoteMedia)

	if !conf.Dryrun {
		uploadMediaItems(newMedia)
		writeRemoteMedia(newMedia)
		for _, m := range newMedia {
			log.Info(m.LocalFile, m.Link)
		}
	}

}
