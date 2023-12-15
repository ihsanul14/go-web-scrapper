FROM golang:alpine as builder

COPY . $GOPATH/src/go-web-scrapper
WORKDIR $GOPATH/src/go-web-scrapper

COPY . .
RUN go mod tidy
RUN apk add --no-cache build-base

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /go-web-scrapper .

FROM alpine:3.4

COPY --from=builder /go-web-scrapper /go-web-scrapper
COPY .env .env
RUN touch .env

EXPOSE 30001

ENTRYPOINT ["./go-web-scrapper"]