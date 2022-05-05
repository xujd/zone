# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN go build -buildvcs=false -o /zone

EXPOSE 1323
CMD [ "/zone","--dbhost=192.168.144.1" ]