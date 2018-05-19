# start a golang base image, version 1.8
FROM golang:latest

#switch to our app directory
RUN mkdir -p /go/src/imageService
WORKDIR /go/src/imageService

#copy the source files
COPY Service.go /go/src/imageService

#disable crosscompiling
ENV CGO_ENABLED=0

#compile linux only		
ENV GOOS=linux

#get go dependencies
RUN go get github.com/nfnt/resize

#build the binary with debug information removed
RUN go build  -ldflags '-w -s' -a -installsuffix cgo -o ImageService