# Get Atom format for github-blog-open-source RSS:
GET http://localhost:8080/atom/github-blog-open-source?secret=1234

HTTP/1.1 200
[Asserts]
status == 200
header "Content-Type" contains "application/atom+xml; charset=utf-8"
body contains "<entry>"

# Get Atom format for github-blog-engineering:
GET http://localhost:8080/atom/github-blog-engineering?secret=1234

HTTP/1.1 200
[Asserts]
status == 200
header "Content-Type" contains "application/atom+xml; charset=utf-8"
body contains "<entry>"
