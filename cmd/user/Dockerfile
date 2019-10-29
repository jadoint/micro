FROM golang:latest
WORKDIR /home/app
COPY . .
RUN go build -o micro_user_service /home/app/cmd/user/main.go
CMD ["./micro_user_service"]
EXPOSE 8000