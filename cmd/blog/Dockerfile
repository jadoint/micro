FROM golang:latest
WORKDIR /home/app
COPY . .
RUN go build -o micro_blog_service /home/app/cmd/blog/main.go
CMD ["./micro_blog_service"]
EXPOSE 8001