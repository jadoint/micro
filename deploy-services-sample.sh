#!/bin/bash

# Build Go binaries
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/blog/blog_service cmd/blog/main.go
env GOOS=linux GOARCH=amd64 go build -o deployments/bin/user/user_service cmd/user/main.go

# Upload to server
ssh host 'mkdir -p deployments/{bin,new_bin,cache,tls}'
scp -r deployments/bin/* youruser@host:~/deployments/new_bin &
scp -r deployments/cache/* youruser@host:~/deployments/cache &
scp -r deployments/tls/* youruser@host:~/deployments/tls &
wait

# Swap old bin and new bin directories
# then delete old bin (to minimize downtime).
ssh host 'cd deployments;mv bin old_bin;mv new_bin bin;rm -rf old_bin'

# Allow a regular user to restart these services without a password
# by allowing it in the sudoers file:
# > visudo
# Add these lines at the end
# youruser    ALL=NOPASSWD: /bin/systemctl restart blog.service
# youruser    ALL=NOPASSWD: /bin/systemctl restart user.service
ssh host 'sudo systemctl restart blog.service;sudo systemctl restart user.service'

echo "All done"
