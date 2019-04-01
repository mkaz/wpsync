
# wpsync

A command-line tool to sync a local directory to your WordPress.

This can publish a post, or upload media.

## Setup

Works with any self hosted WordPress but requires the [JWT Authentication](https://wordpress.org/plugins/jwt-authentication-for-wp-rest-api/) plugin to be installed and activated. Follow the plugin instructions for installation and setup.

Configure wpsync to work with you site using: `wpsync --init` It will prompt you for your username and password, the password is not stored but the JWT token used to make API calls. The token expires after 7 days, so you will need to login again.

Create a `media` sub-directory, anything placed in here will be copied to the media library.

Create a `posts` sub-directory, each markdown file placed here will create a new post.

The markdown files accept a front matter to specify settings. The front matter format is similar to Jekyll, a set of parameters delineated by lines containing `---`

The parameters are: `title, category, date, tags, status`

See [WordPress REST API](https://developer.wordpress.org/rest-api/reference/posts/#create-a-post) for parameter details and default values.

Run `wpsync`

The program will create a `posts.json` and `media.json` file locally with the entries that were uploaded. If these json files are deleted, the files found in posts & media directories will be uploaded again.


###

You can confirm the JWT Authentication plugin is installed and working properly, by using this curl command and checking to see if you get a proper token response, replace USER/PASS with your credentials.

```
curl -X POST -d "username=USER&password=PASS" http://your.site/wp-json/jwt-auth/v1/token
```


## Colophon

* Created by Marcus Kazmierczak ([mkaz.blog](https://mkaz.blog/))
* Written in Golang.
* Pull requests & Bug reports welcome.
* WTFPL Licensed.

