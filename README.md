<img class="logo" src="logo.svg" width="355px" height="142px" alt="FeedPLS"/>

Feedpls is a simple web service to generate RSS/Atom for web page which does not have it.

[![Build Status](https://github.com/witjem/feedpls/actions/workflows/build/badge.svg)](https://github.com/witjem/feedpls/actions)
[![Image Size](https://img.shields.io/docker/image-size/witjem/feedpls/main)](https://hub.docker.com/r/witjem/feedpls)
## Configuration

Example `feeds.yaml`

```yaml
---
- id: "" # Feed identifier.
  title: "" # The title in the feed.
  description: "" # The description in the feed.
  url: "" # The URL of the web page to generate the feed from.
  
  matchers:
  
    # Defines how to find items URLs.
    itemUrl: 
      selector: "main a" # A CSS selector.
      attr: "href" # Optional. HTML attribute name. Use it to get content from attribute.

    # Defines how to find item title. 
    # The Title value is retrieved from the Item page.
    title:
      selector: "main .c-body-title h1"
      attr: ""

    # Defines how to find item description.
    # The Description value is retrieved from the Item page.
    description:
      selector: "meta[name='twitter:description']"
      attr: "content"

    # Defines how to find item published time.
    # The Published value is retrieved from the Item page.
    published:
      selector: "meta[name='article:published_time']"
      attr: "content"
      
      # Standard golang layout for parsing time. 
      # Examples: https://pkg.go.dev/time#pkg-constants
      layout: "2006-01-02T15:04:05Z07:00" 

```

## Run

Example `docker-compose.yml` 

```yml
services:
  feedpls:
    image: witjem/feedpls:main
    container_name: feedpls
    hostname: feedpls
    ports:
      - "8080:8080"
    volumes:
      - ./feeds.yaml:/feeds.yaml
    environment:
      - APP_PORT=8080
      - APP_SECRET=1234 # access-key
      - APP_FEEDS=/feeds.yaml
      - APP_TTL=5m # feed cache time
```

## Endpoints
```shell
## Get RSS feed
curl -X GET --location "https://<host>/rss/<feed-id>?secret=<access-key>"

## Get Atom feed
curl -X GET --location "https://<host>/atom/<feed-id>?secret=<access-key>"
```

## Application options
```
--feeds value   yaml file with describes feeds matching (default: "./feeds.yml") [$APP_FEEDS]
--help, -h      show help (default: false)
--port value    server port (default: "8080") [$APP_PORT]
--secret value  api secret (default: "secret") [$APP_SECRET]
--ttl value     feed caching time (default: 5m0s) [$APP_TTL]
```
## Alternatives
* [RSS Please](https://github.com/wezm/rsspls)
