# Build our application using a Go builder.
FROM golang:1.19 AS builder
WORKDIR /src/litefs-example
COPY . .
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/litefs-example ./cmd/litefs-example


# Our final Docker image stage starts here.
FROM alpine

# Copy binaries from the previous build stages.
COPY --from=flyio/litefs:0.3 /usr/local/bin/litefs /usr/local/bin/litefs
COPY --from=builder /usr/local/bin/litefs-example /usr/local/bin/litefs-example

# Copy our LiteFS configuration.
ADD etc/litefs.yml /etc/litefs.yml

# Setup our environment to include FUSE & SQLite. We install ca-certificates
# so we can communicate with the Consul server over HTTPS. cURL is added so
# we can call our HTTP endpoints for debugging.
RUN apk add bash fuse sqlite ca-certificates curl

# Run LiteFS as the entrypoint. Anything after the double-dash is run as a
# subprocess by LiteFS. This allows the file system to mount and initialize
# before the application starts.
ENTRYPOINT litefs mount -- litefs-example -dsn /litefs/db
