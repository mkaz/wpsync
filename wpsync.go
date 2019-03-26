// wpsync - command-line tool to sync wordpress
// https://github.com/mkaz/wpsync
//
// TODO: add watch
// TODO: add confirmation flag
// TODO: add verbose flag

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BlogID string `toml:"blog_id"`
	Token  string `toml:"token"`
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
		log.Fatalln("No auth token configured in wpsync.conf")
	}

	if conf.BlogID == "" {
		log.Fatalln("No blog id configured in wpsync.conf")
	}
}

// route command and args
func main() {

	// read local files for data
	localPosts := getLocalPosts()
	remotePosts := getRemotePosts()
	newPosts := comparePosts(localPosts, remotePosts)
	uploadPosts(newPosts)
	writeRemotePosts(newPosts)

	// media
	localMedia := getLocalMedia()
	remoteMedia := getRemoteMedia()
	newMedia := compareMedia(localMedia, remoteMedia)
	uploadMediaItems(newMedia)
	writeRemoteMedia(newMedia)
	for _, m := range newMedia {
		fmt.Println(m.LocalFile, m.Link)
	}

}
