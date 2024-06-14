FROM golang:1.22-alpine as builder

ENV GOPROXY="https://goproxy.cn,direct"
ARG VERSION
WORKDIR /app/
COPY . .
RUN go generate -tags plugins ./engine/plugins
RUN go mod download

RUN CGO_ENABLED=0 go build -o ./bin/wxbot

FROM alpine:latest
LABEL MAINTAINER="nibuchiwochile@gmail.com"
WORKDIR /app/
VOLUME ["/app/data/", "/app/config"]

RUN apk add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/shanghai" >> /etc/timezone \
    && apk del tzdata

COPY --from=builder /app/bin/wxbot /app/wxbot

EXPOSE 9528
ENTRYPOINT ["./wxbot"]