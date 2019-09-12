#!/bin/bash

# Build Go binaries
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/blog/blog_service cmd/blog/main.go
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/user/user_service cmd/user/main.go

# Upload to server
ssh host 'mkdir -p deployments'
scp -r deployments/bin user@host:~/deployments/bin &
scp -r deployments/cache user@host:~/deployments/cache &
wait

echo "All done"
