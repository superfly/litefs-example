# Build our application using a Go builder.
FROM golang:1.20 AS builder

WORKDIR /src/litefs-example
COPY . .
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/litefs-example ./cmd/litefs-example


# Our final Docker image stage starts here.
FROM alpine
ARG LITEFS_CONFIG=litefs.yml

# Copy binaries from the previous build stages.
COPY --from=flyio/litefs:0.4 /usr/local/bin/litefs /usr/local/bin/litefs
COPY --from=builder /usr/local/bin/litefs-example /usr/local/bin/litefs-example

# Copy our LiteFS configuration.
ADD fly-io-config/etc/litefs.yml /tmp/litefs.yml
ADD docker-config/etc/litefs.primary.yml /tmp/litefs.primary.yml
ADD docker-config/etc/litefs.replica.yml /tmp/litefs.replica.yml

RUN cp /tmp/$LITEFS_CONFIG /etc/litefs.yml

# Setup our environment to include FUSE & SQLite. We install ca-certificates
# so we can communicate with the Consul server over HTTPS. cURL is added so
# we can call our HTTP endpoints for debugging.
RUN apk add bash fuse3 sqlite ca-certificates curl

# Run LiteFS as the entrypoint. After it has connected and sync'd with the
# cluster, it will run the commands listed in the "exec" field of the config.
ENTRYPOINT litefs mount
