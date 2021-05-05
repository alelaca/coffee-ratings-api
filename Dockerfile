FROM golang

ADD . /go/src/github.com/alelaca/coffee-ratings-api

RUN go install github.com/alelaca/coffee-ratings-api/cmd@main

ENTRYPOINT /go/bin

EXPOSE 9000
