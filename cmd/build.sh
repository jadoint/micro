#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -o blog/blog_service blog/main.go
env GOOS=linux GOARCH=amd64 go build -o user/user_service user/main.go