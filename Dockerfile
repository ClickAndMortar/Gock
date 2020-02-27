FROM golang:1.13 as go-builder

WORKDIR /go/src/app

COPY main.go ./

RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:buster-slim

COPY --from=go-builder /go/bin/app /usr/local/bin/gock

EXPOSE 8000

CMD ["gock"]
