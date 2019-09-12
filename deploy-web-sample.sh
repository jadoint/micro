#!/bin/bash

# Upload to server
ssh host 'rm -rf web/public;rm -rf web/src;rm web/package.json'
scp -r web/public user@host:~/web/public &
scp -r web/src user@host:~/web/src &
scp -r web/package.json user@host:~/web/package.json &
wait

# Install node modules
ssh host 'cd web;yarn install'

# Leave build directory last to minimize site downtime
ssh host 'rm -rf web/build'
ssh host 'cd web;yarn build'

echo "All done"
