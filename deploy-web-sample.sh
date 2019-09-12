#!/bin/bash

# Build Go binaries
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/blog/blog_service cmd/blog/main.go
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/user/user_service cmd/user/main.go

# Upload to server
ssh jadoweb 'mkdir -p deployments;mkdir -p web'
scp -r deployments/bin user@host:~/deployments/bin &
scp -r deployments/cache user@host:~/deployments/cache &
scp -r web/public user@host:~/web/public &
scp -r web/src user@host:~/web/src &
scp -r web/package.json user@host:~/web/package.json &
wait

ssh host 'cd web;yarn install;yarn build'

echo "All done"
