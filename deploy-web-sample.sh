#!/bin/bash

# Create production build
cd web
yarn build

# Upload to server
ssh host 'mkdir -p web/{build,new_build}'
scp -r build/* youruser@host:~/web/new_build

# Swap old build and new build directories
# then delete old build (to minimize downtime).
ssh host 'cd web;mv build old_build;mv new_build build;rm -rf old_build'

echo "All done"
