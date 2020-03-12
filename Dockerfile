FROM golang:latest
MAINTAINER Raghavendra Balgi;rkbalgi@gmail.com
ADD html /etc/isosim/web
ADD specs /etc/isosim/specs
ADD testdata /etc/isosim/data
COPY . /home/isosim/app/src/github.com/rkbalgi/isosim
ENV GOPATH /home/isosim/app
ENV HTTP_PORT 8080
# go get other dependencies
RUN go get github.com/rkbalgi/go/net && go get github.com/rkbalgi/go/hsm \
&& go get github.com/rkbalgi/go/encoding/ebcdic && go get github.com/sirupsen/logrus
WORKDIR /home/isosim/app/src/github.com/rkbalgi/isosim/cmd/isosim
ENTRYPOINT go run isosim.go -httpPort $HTTP_PORT -htmlDir /etc/isosim/web -dataDir /etc/isosim/data -specDefFile /etc/isosim/specs/isoSpecs.spec
