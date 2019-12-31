# Builder
FROM golang:latest AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o micro_user_service /build/cmd/user/main.go

# Production build
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /home/app
COPY --from=builder /build/micro_user_service .
EXPOSE 9000
ENTRYPOINT ["./micro_user_service"]