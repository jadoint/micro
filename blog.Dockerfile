FROM golang:alpine AS builder
WORKDIR /home/app
COPY . .

# Create self-signed certs below if no longer using volumes for SSL
RUN apk add --update openssl && \
  rm -rf /var/cache/apk/*
RUN mkdir -p /home/ssl
RUN openssl req -new -newkey rsa:2048 -x509 -days 3650 -nodes -out /home/ssl/site.pem -keyout /home/ssl/site.key -subj "/C=US/ST=Wyoming/L=Cheyenne/O=Jado Interactive/OU=IT/CN=asianfanfics.com"

# build-base includes gcc which kafka needs for builds
RUN apk add build-base

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# "-tags musl" needed for kafka when doing alpine builds
RUN GOOS=linux go build -tags musl -a -o micro_blog_service /home/app/cmd/blog/main.go

# Make new container from builder
FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /home/app

RUN mkdir -p /home/app/logs

RUN mkdir -p /home/ssl
COPY --from=builder /home/ssl /home/ssl

COPY --from=builder /home/app/micro_blog_service /home/app

EXPOSE 9000
ENTRYPOINT ./micro_blog_service