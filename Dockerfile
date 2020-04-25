# stage 1
FROM golang:alpine
MAINTAINER Raghavendra Balgi;rkbalgi@gmail.com
COPY . /home/isosim/app/isosim
WORKDIR /home/isosim/app/isosim/cmd/isosim
RUN go build -v -o app .

# stage 2
FROM alpine
#USER 1001:1001
ADD web /etc/isosim/web
ADD test/testdata/specs /etc/isosim/specs
ADD test/testdata/appdata /etc/isosim/data
ADD test/testdata/certs /etc/isosim/certs
ENV HTTP_PORT 8080
ENV LOG_LEVEL DEBUG
ENV TLS_ENABLED=false
ENV TLS_CERT_FILE=/etc/isosim/certs/cert.pem
ENV TLS_KEY_FILE=/etc/isosim/certs/key.pem
COPY --from=0 /home/isosim/app/isosim/cmd/isosim/app /usr/apps/isosim
ENTRYPOINT /usr/apps/isosim --http-port $HTTP_PORT \
                            --log-level $LOG_LEVEL \
                            --html-dir /etc/isosim/web \
                            --data-dir /etc/isosim/data \
                            --specs-dir /etc/isosim/specs
