FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/tidwall/buntdb
ADD ./src /app
RUN go build ./app.go

CMD ["./app"]
