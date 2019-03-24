
# wpsync

A command-line tool to sync a local directory to your WordPress.

This can publish a post, or upload media.

## Setup

Works with WordPress.com and self-hosted blogs running [Jetpack](https://jetpack.me/). If running Jetpack it requires the JSON API to be enabled, which should be activated by default.

TODO: Improve authentication. For now you need to configure with an authorization token and blog id for each of your blogs. To do so, you can authorize using this site run by Apokalyptik, its fine trust me, which you must if you are going to run my code: https://rest-access.test.apokalyptik.com/

Once you obtain the token and blog id, create a directory for the site you want to sync and add the file `wpsync.conf` with the following two parameters:

    token = ABCDEFGH123456
    blog_id = 123456


Create a `media` sub-directory, anything placed in here will be copied to the media library.

Create a `posts` sub-directory, each markdown file placed here will create a new post.

Run `wpsync` to do it.


## Colophon

* Created by Marcus Kazmierczak
* Written in Golang.
* Pull requests & Bug reports welcome.
* WTFPL Licensed.

