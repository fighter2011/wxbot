FROM golang:1.22-alpine as builder

ENV GOPROXY="https://goproxy.cn,direct"
ARG VERSION
WORKDIR /app/

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go generate -tags plugins ./engine/plugins

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/wxbot

FROM alpine:latest
LABEL MAINTAINER="nibuchiwochile@gmail.com"
WORKDIR /app/
VOLUME ["/app/data/", "/app/config"]

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/shanghai" >> /etc/timezone \
    && apk del tzdata \
    && apk add mailcap

COPY --from=builder /app/bin/wxbot /app/wxbot

EXPOSE 9528
ENTRYPOINT ["./wxbot"]