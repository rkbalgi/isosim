FROM golang:latest
MAINTAINER Raghavendra Balgi;rkbalgi@gmail.com
ADD html /etc/isosim/web
ADD specs /etc/isosim/specs
COPY . /home/isosim/app/src/github.com/rkbalgi/isosim
ENV GOPATH /home/isosim/app
ENV HTTP_PORT 8080
# go get other dependencies
RUN go get github.com/rkbalgi/go/net
RUN go get github.com/rkbalgi/go/hsm
RUN go get github.com/rkbalgi/go/encoding/ebcdic
WORKDIR /home/isosim/app/src/github.com/rkbalgi/isosim
ENTRYPOINT go run isosim.go -httpPort $HTTP_PORT -htmlDir /etc/isosim/web -dataDir /home/isosim -specDefFile /etc/isosim/specs/isoSpecs.spec
