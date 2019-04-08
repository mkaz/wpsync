
# wpsync

A command-line tool to sync a local directory to your WordPress.

This can publish and update posts, as well as upload media.

## Install

Go Environment: install using: `go get github.com/mkaz/wpsync`

Binaries: See [releases tab](https://github.com/mkaz/wpsync/releases) in Github


## Setup

Works with any self-hosted WordPress but requires the [JWT Authentication](https://wordpress.org/plugins/jwt-authentication-for-wp-rest-api/) plugin to be installed and activated. Follow the plugin instructions for installation and setup.

### Plugin

Here's what I did to install the JWT Auth plugin:
```
wp plugin install --activate jwt-authentication-for-wp-rest-api
```

Add the following to your `wp-config.php`

```
define('JWT_AUTH_SECRET_KEY', 'your-top-secret-key');
```


My `.htaccess` config to configure the plugin:
```
# BEGIN WordPress
<IfModule mod_rewrite.c>
RewriteEngine On
RewriteBase /

RewriteCond %{HTTP:Authorization} ^(.*)
RewriteRule ^(.*) - [E=HTTP_AUTHORIZATION:%1]

RewriteRule ^index\.php$ - [L]
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule . /index.php [L]
</IfModule>
# END WordPress
```

### wpsync Setup

Configure wpsync to work with you site using: `wpsync --init` It will prompt you for your username and password, the password is not stored but the JWT token used to make API calls. The token expires after 7 days, so you will need to login again.

Create a `media` sub-directory, anything placed in here will be copied to the media library.

Create a `posts` sub-directory, each markdown file placed here will create a new post.


## Usage

Run `wpsync [args]`

Arguments:

	-confirm
		Confirm prompt before upload
	-dryrun
		Test run, shows what will happen
	-help
		Display help and quit
	-init
		Create settings for blog and auth
	-verbose
		Details lots of details
	-version
		Display version and quit

### Posts Markdown

The posts should be written in markdown and include "front-matter" to specify settings. The front-matter format is similar to Jekyll, a set of parameters delineated by lines containing `---`

The parameters are: `title, date, status`

See [WordPress REST API](https://developer.wordpress.org/rest-api/reference/posts/#create-a-post) for parameter details and default values.

For example, if you want to publish a draft set `status: draft` in the front-matter in the markdown. Edit, and preview away, and then when ready to publish, change to `status: publish`.

Post example:

```
---
title: My Sample Post
status: draft
---

Content for my post...
```

### Sync Data

The program creates a `posts.json` and `media.json` files locally with the entries that were uploaded. If these json files are deleted, then any files found in posts & media directories will be uploaded again.

TODO: Implement two-way sync, right now the data only goes from local to remote.

## Troubleshoot

You can confirm the JWT Authentication plugin is installed and working properly, by using this curl command and checking to see if you get a proper token response, replace USER/PASS with your credentials.

```
curl -X POST -d "username=USER&password=PASS" http://your.site/wp-json/jwt-auth/v1/token
```


## Colophon

* Created by Marcus Kazmierczak ([mkaz.blog](https://mkaz.blog/))
* Written in Golang.
* Pull requests & Bug reports welcome.
* WTFPL Licensed.

