# stage 1
FROM golang:alpine
MAINTAINER Raghavendra Balgi;rkbalgi@gmail.com
COPY . /home/isosim/app/isosim
WORKDIR /home/isosim/app/isosim/cmd/isosim
RUN go build -v -o app isosim.go

# stage 2
FROM alpine
#USER 1001:1001
COPY --from=0 /home/isosim/app/isosim/cmd/isosim/app /usr/apps/isosim
ADD html /etc/isosim/web
ADD specs /etc/isosim/specs
ADD testdata /etc/isosim/data
ADD certs /etc/isosim/certs
ENV HTTP_PORT 8080
# ENV TLS_ENABLED=true
# ENV TLS_CERT_FILE=/etc/isosim/certs/cert.pem
# ENV TLS_KEY_FILE=/etc/isosim/certs/key.pem
ENTRYPOINT /usr/apps/isosim -http-port $HTTP_PORT -html-dir /etc/isosim/web -data-dir /etc/isosim/data -specs-dir /etc/isosim/specs
