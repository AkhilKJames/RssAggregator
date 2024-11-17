# RSS AGGREGATOR

## What is it?

An RSS feed aggregator in Go! It's a web server that allows clients to:

* Add RSS feeds to be collected
* Follow and unfollow RSS feeds that other users have added
* Fetch all of the latest posts from the RSS feeds they follow
* RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## Setup

Before we dive into the project, let's make sure you have everything you'll need on your machine.

1. An editor. I use [VS code](https://code.visualstudio.com/), you can use whatever you like.
2. The latest [Go toolchain](https://golang.org/doc/install).
3. If you're in VS Code, I recommend the official [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go).
4. An HTTP client. I use [POSTMAN](https://www.postman.com/downloads/), but you can use whatever you like.
5. [PostgreSQL and few database dependencies](./docs/postgres.md)


## Run and test server

```bash
go build -o out && ./out
```
Alternatively -
```bash
You can directly run in debug mode with the launch.json
```

Once it's running, use an HTTP client to test your endpoints.

## What it does?

* You can a new user.
* A user can add feeds/feed urls into the system using his ApiKey which he gets on account creation.
* A user can follow already existing feeds using the feed id
* Server will continously fetch the latest feeds using the url provided by users.
* Users can get the posts they are interested in by following the feeds they like.
* [API Doc](./docs/api.md) for more details :)
