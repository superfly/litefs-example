# Fetch the LiteFS binary using a multi-stage build.
FROM flyio/litefs:pr-94 AS litefs


# Build our "myapp" application using a Go builder.
FROM golang:1.19 AS builder
WORKDIR /src/myapp
COPY . .
RUN go build -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/myapp ./cmd/myapp


# Our final Docker image stage starts here.
FROM alpine

# Copy binaries from the previous build stages.
COPY --from=builder /usr/local/bin/myapp /usr/local/bin/myapp
COPY --from=litefs /usr/local/bin/litefs /usr/local/bin/litefs

# Copy our LiteFS configuration.
ADD etc/litefs.yml /etc/litefs.yml

# Setup our environment to include FUSE & SQLite.
RUN apk add bash curl fuse sqlite

# Ensure our data directory exists before mounting with LiteFS.
RUN mkdir -p /data

# Run LiteFS as the entrypoint so it can execute "myapp" as a subprocess.
ENTRYPOINT "litefs"
