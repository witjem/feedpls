# Get RSS format for github-blog-open-source:
GET http://localhost:8080/rss/github-blog-open-source?secret=1234

HTTP/1.1 200
[Asserts]
status == 200
header "Content-Type" contains "application/rss+xml; charset=utf-8"
body contains "<item>"

# Get RSS format for github-blog-engineering:
GET http://localhost:8080/rss/github-blog-engineering?secret=1234

HTTP/1.1 200
[Asserts]
status == 200
header "Content-Type" contains "application/rss+xml; charset=utf-8"
body contains "<item>"
