# Step 1: Modules caching
FROM golang:1.19.1-alpine3.16 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.19.1-alpine3.16 as builder

ARG GIT_BRANCH
ARG GITHUB_SHA

ENV CGO_ENABLED=0

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN \
    version=${GIT_BRANCH}-${GITHUB_SHA:0:7}-$(date +%Y%m%dT%H:%M:%S); \
    echo "version=$version" && \
    go build -o /bin/app -ldflags "-X main.revision=${version} -s -w" ./cmd/server

# Step 3: Final
FROM scratch
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]
