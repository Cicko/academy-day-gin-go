FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

ADD ./src /app
RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/tidwall/buntdb
RUN go build ./app.go

EXPOSE 8080

CMD ["./app"]
