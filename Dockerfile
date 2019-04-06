FROM golang:1.12.2-alpine3.9

WORKDIR /go/src/app
COPY . .

RUN apk update && apk add git gcc libc-dev
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]