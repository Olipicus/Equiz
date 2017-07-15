FROM alpine

RUN apk add --no-cache ca-certificates && apk update && rm -rf /var/cache/apk/*
ADD app /app
ENTRYPOINT ["/app"]