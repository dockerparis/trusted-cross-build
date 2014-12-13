FROM golang:1.3-cross

ENV CGO_ENABLED 0

ADD wrapper.go /go/

# gosu-$(dpkg --print-architecture)
RUN GOARCH=amd64       go build -v -ldflags -d -o /go/bin/wrapper-amd64
RUN GOARCH=386         go build -v -ldflags -d -o /go/bin/wrapper-i386
RUN GOARCH=arm GOARM=5 go build -v -ldflags -d -o /go/bin/wrapper-armel
RUN GOARCH=arm GOARM=6 go build -v -ldflags -d -o /go/bin/wrapper-armhf
