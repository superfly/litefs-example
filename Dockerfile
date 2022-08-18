FROM flyio/litefs:pr-63 as litefs


FROM alpine
COPY --from=litefs /usr/local/bin/litefs /usr/local/bin/litefs

RUN apk add bash curl fuse sqlite
RUN mkdir -p /data
COPY etc/litefs.yml /etc/litefs.yml
ENTRYPOINT [ "litefs" ]
