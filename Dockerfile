FROM golang:1.15 as build
WORKDIR /go/src/github.com/linuxdynasty/hello-world-service/
RUN go get github.com/netdata/statsd
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hello-world .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build go/src/github.com/linuxdynasty/hello-world-service/hello-world .
EXPOSE 80
CMD ["./hello-world"]  