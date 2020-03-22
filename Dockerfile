FROM golang:latest
MAINTAINER Raghavendra Balgi;rkbalgi@gmail.com
ADD html /etc/isosim/web
ADD specs /etc/isosim/specs
ADD testdata /etc/isosim/data
COPY . /home/isosim/app/isosim
ENV HTTP_PORT 8080
# go get other dependencies
WORKDIR /home/isosim/app/isosim/cmd/isosim
ENTRYPOINT go run isosim.go -httpPort $HTTP_PORT -htmlDir /etc/isosim/web -dataDir /etc/isosim/data -specDefFile /etc/isosim/specs/isoSpecs.spec
