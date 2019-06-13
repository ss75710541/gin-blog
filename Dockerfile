FROM golang:1.12

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

WORKDIR /go/cache
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ /go/src/gin-blog
WORKDIR /go/src/gin-blog

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o gin-blog

FROM alpine:3.9

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/gin-blog/gin-blog /
COPY --from=0 /go/src/gin-blog/conf /conf

CMD ["/gin-blog"]
