// wpsync - command-line tool to sync wordpress
// https://github.com/mkaz/wpsync
//
// TODO: add watch
// TODO: add confirmation flag

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/automattic/go/jaguar"
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
	Id        int       `json:"ID"`
	Date      time.Time `json:"date"`
	URL       string    `json:"URL"`
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

	fmt.Println("Blog ID: " + conf.BlogID)

	// check posts
	localPosts := getLocalPosts()
	remotePosts := getRemotePosts()

	// check media
	localMedia := getLocalMedia()

	fmt.Printf("%v \n", localPosts)
	fmt.Printf("%v \n", remotePosts)

	// - read posts.json
	// - compare directory to json
	//		- check against local info
	//		- check against remote info
	//		- sync (upload or download)

	// check media
	//	- see if any media files
	//  - check against local info
	//  - copy new files up
}

// read posts from local directory
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

// read media from local directory
func getLocalMeida() (media []Media) {
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

// read posts from json file
func getRemotePosts() (posts []Post) {
	return posts
}

func getApiFetcher(endpoint string) (j jaguar.Jaguar) {
	apiurl := "https://public-api.wordpress.com/rest/v1"
	url := strings.Join([]string{apiurl, "sites", conf.BlogID, endpoint}, "/")

	j = jaguar.New()
	j.Header.Add("Authorization", "Bearer "+conf.Token)
	j.Url(url)
	return j
}
