FROM golang:latest
ADD html /etc/isosim/web
ADD specs /etc/isosim/specs
COPY . /home/isosim/app/src/github.com/rkbalgi/isosim
#COPY ../go /home/isosim/app/src/github.com/rkbalgi/go
ENV GOPATH /home/isosim/app
ENV HTTP_PORT 8080
RUN go get github.com/rkbalgi/go/net
RUN go get github.com/rkbalgi/go/hsm
RUN go get github.com/rkbalgi/go/encoding/ebcdic
WORKDIR /home/isosim/app/src/github.com/rkbalgi/isosim
RUN ls
ENTRYPOINT go run isosim.go -httpPort $HTTP_PORT -htmlDir /etc/isosim/web -dataDir /home/isosim -specDefFile specs/isoSpecs.spec
