FROM alpine:3.5

RUN apk add -U --no-cache ca-certificates

COPY cryptocompare-scrape /cryptocompare-scrape

CMD ["/cryptocompare-scrape"]    