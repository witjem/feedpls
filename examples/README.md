# Example

This example shows how made feed for https://github.blog

## Run

```shell
docker compose up -d
```

Feed for https://github.blog/category/open-source/:

* RSS http://localhost:8080/rss/github-blog-open-source?secret=1234
* Atom http://localhost:8080/atom/github-blog-open-source?secret=1234

Feed for https://github.blog/category/open-source/:

* RSS http://localhost:8080/rss/github-blog-engineering?secret=1234
* Atom http://localhost:8080/atom/github-blog-engineering?secret=1234
