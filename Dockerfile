FROM golang:1.12.7 AS builder
WORKDIR $GOPATH/src/github.com/kuhnuri/batch-graphics
RUN go get -v -u github.com/kuhnuri/go-worker
COPY docker/main.go .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN go build -a -o main .

#FROM adoptopenjdk/openjdk11:jdk-11.0.3_7-slim
FROM ubuntu:18.04

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'
RUN apt-get -y update \
    && apt-get -y install --no-install-recommends ca-certificates locales imagemagick \
    && apt-get -y clean \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen en_US.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/app
COPY --from=builder /go/src/github.com/kuhnuri/batch-graphics/main .

ENTRYPOINT ["./main"]
