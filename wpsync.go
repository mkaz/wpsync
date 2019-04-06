// wpsync - command-line tool to sync wordpress
// https://github.com/mkaz/wpsync
//

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// Config is the structure of the jwt-auth response and
// settings, it is used to unmarshal the data
type Config struct {
	SiteURL string `json:"site-url"`
	Token   string `json:"token"`
}

type Post struct {
	Id        int    `json:"id"`
	Title     string `json:"title.raw"`
	Date      WPTime `json:"date_gmt"`
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
var dryrun bool
var confirm bool

// read config and parse args
func init() {

	var helpFlag = flag.Bool("help", false, "Display help and quit")
	var versionFlag = flag.Bool("version", false, "Display version and quit")
	var testFlag = flag.Bool("test", false, "Test config and authentication")
	flag.BoolVar(&log.Verbose, "verbose", false, "Details lots of details")
	flag.BoolVar(&dryrun, "dryrun", false, "Test run, shows what will happen")
	flag.BoolVar(&setup, "init", false, "Create settings for blog and auth")
	flag.BoolVar(&confirm, "confirm", false, "Confirm prompt before upload")
	flag.Parse()

	if *helpFlag {
		usage()
	}

	if *versionFlag {
		fmt.Println("wpsync v0.1.0")
		os.Exit(0)
	}

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

	if *testFlag {
		if testSetup() {
			fmt.Println("Test setup passed. 👍")
			os.Exit(0)
		} else {
			fmt.Println("Test setup fail. 👎")
			os.Exit(1)
		}
	}

	if setup {
		runSetup()
	}

	// test setup
	if !testSetup() {
		// setup not working
		// check if runSetup() ran with setup
		// if not run it now otherwise bail
		log.Fatal("Error validating.", conf)
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

	if !dryrun {
		newPosts = createPosts(newPosts)
		updatedPosts = updatePosts(updatedPosts)
		writeRemotePosts(newPosts, updatedPosts)
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

	if !dryrun {
		uploadedMedia := uploadMediaItems(newMedia)
		writeRemoteMedia(uploadedMedia)
	}
}

func confirmPrompt(prompt string) bool {

	// confirmation not required
	if !confirm {
		return true
	}

	var ans string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&ans)
	if err != nil {
		log.Fatal("What happened?", err)
	}
	if ans == "y" || ans == "Y" {
		return true
	} else {
		return false
	}
}

// Display Usage
func usage() {
	fmt.Println("usage: wpsync [args] \n")
	fmt.Println("Arguments:\n")
	flag.PrintDefaults()
	fmt.Println("")
	os.Exit(0)
}
