<img class="logo" src="logo.svg" width="355px" height="142px" alt="FeedPLS"/>

Feedpls is a simple web service to generate RSS/Atom for web page which does not have it.

[![Build Status](https://github.com/witjem/feedpls/actions/workflows/ci.yml/badge.svg)](https://github.com/witjem/feedpls/actions)
[![Image Size](https://img.shields.io/docker/image-size/witjem/feedpls/main)](https://hub.docker.com/r/witjem/feedpls)

## Configuration

Example `feeds.yml`

```yaml
---
# Define config with `XPath` parser engine. 
- id: "news-first" # Feed identifier.
  title: "" # The title in the feed.
  description: "" # The description in the feed.
  url: "" # The URL of the web page to generate the feed from.

  matchers:
    engine: "xpath"
    
    # Defines how to find items URLs.
    itemUrl:
      expr: "//main//a/@href" # XPath expression.

    # Defines how to find item title. 
    # The Title value is retrieved from the Item page.
    title:
      expr: "//title"

    # Defines how to find item description.
    # The Description value is retrieved from the Item page.
    description:
      expr: "//meta[@name='twitter:description']/@content"

    # Defines how to find item published time.
    # The Published value is retrieved from the Item page.
    published:
      expr: "//meta[@name='article:published_time']/@content"

      # Standard golang layout for parsing time. 
      # Examples: https://pkg.go.dev/time#pkg-constants
      layout: "2006-01-02T15:04:05Z07:00"
      tz: "Europe/Kyiv" # Optional, default value UTC
      
  # Optional property. Functions to be applied to the matched data after parsing.
  # For example if after parse you want to replace (or remove) some text to another.
  functions: 
    - replace: # Replace function which replaces the matched data with the given value.
        field: "title" # Defines which property to apply to. Can be 'title' or 'description'
        from: " -- ABC News" # Defines what text want to replace. 
        to: "" # Defines what text want to replace to. 
  
# Define config with `GoQuery` (JQuery like) parser engine. 
- id: "news-second" # Feed identifier.
  title: "" # The title in the feed.
  description: "" # The description in the feed.
  url: "" # The URL of the web page to generate the feed from.
  matchers:
    engine: "goquery"
  
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
      tz: "Europe/Kyiv" # Optional, default value UTC
```

## Run

Example `docker-compose.yml` 

```yml
services:
  feedpls:
    image: witjem/feedpls:v0.2.1 # or ghcr.io/witjem/feedpls:v0.2.1
    container_name: feedpls
    hostname: feedpls
    ports:
      - "8080:8080"
    volumes:
      - ./feeds.yml:/feeds.yml
    environment:
      - APP_PORT=8080
      - APP_SECRET=1234 # access-key
      - APP_FEEDS=/feeds.yml
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

```shell
--feeds value   yaml file with describes feeds matching (default: "./feeds.yml") [$APP_FEEDS]
--help, -h      show help (default: false)
--port value    server port (default: "8080") [$APP_PORT]
--secret value  api secret (default: "secret") [$APP_SECRET]
--ttl value     feed caching time (default: 5m0s) [$APP_TTL]
```

## Alternatives

* [RSS Please](https://github.com/wezm/rsspls)
