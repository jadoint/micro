#!/bin/bash

# Upload to server
ssh host 'rm -rf web;mkdir -p web'
scp -r web/public user@host:~/web/public &
scp -r web/src user@host:~/web/src &
scp -r web/package.json user@host:~/web/package.json &
wait

ssh host 'cd web;yarn install;yarn build'

echo "All done"
