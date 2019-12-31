FROM golang:latest
WORKDIR /home/app

RUN mkdir -p /home/ssl
RUN openssl req -new -newkey rsa:2048 -x509 -days 3650 -nodes -out /home/ssl/site.pem -keyout /home/ssl/site.key -subj "/C=US/ST=Wyoming/L=Cheyenne/O=Jado Interactive/OU=IT/CN=davidado.com"

RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon --build="go build -o micro_blog_service /home/app/cmd/blog/main.go" --command="./micro_blog_service" --exclude-dir=.git --exclude-dir=web --exclude-dir=deployments
EXPOSE 9000