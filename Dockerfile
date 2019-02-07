FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/tidwall/buntdb
RUN go get -u github.com/radovskyb/watcher/...
ADD ./src /app
RUN go build ./app.go

CMD ["./app"]
