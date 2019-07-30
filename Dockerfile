FROM golang:1.12.7 AS builder
WORKDIR $GOPATH/src/github.com/kuhnuri/batch-graphics
RUN go get -v -u github.com/kuhnuri/go-worker
COPY docker/main.go .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN go build -a -o main .

FROM adoptopenjdk/openjdk11:jdk-11.0.3_7-slim
WORKDIR /opt/app
RUN apt-get -y update \
    && apt-get -y install imagemagick \
    && apt-get -y clean \
    && rm -rf /var/lib/apt/lists/*
#RUN convert --help
COPY --from=builder /go/src/github.com/kuhnuri/batch-graphics/main .
#COPY build/dist /opt/app/lib
#COPY docker/run.sh /opt/app/run.sh
#RUN chmod 755 /opt/app/run.sh

ENTRYPOINT ["./main"]
