FROM golang:1.10

LABEL maintainer="Ivan de la Beldad Fernandez <ivandelabeldad@gmail.com>"

ENV GOPATH=/go

ADD . /go/src/sonarr-parser-helper

WORKDIR /go/src/sonarr-parser-helper

RUN go get ./... && \
    go build -o main .

RUN ln -s /tv /televisión && ln -s /movies /películas

CMD ["/go/src/sonarr-parser-helper/main"]
