FROM golang:1.12.2-alpine3.9

WORKDIR /go/src/github.com/kind84/iterpro
COPY . .

RUN apk update && apk add git gcc libc-dev
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["iterpro"]