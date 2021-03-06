#!/bin/bash
#

ssh yourwebserver 'mkdir -p mysql/micro-user/{data,backups};mkdir -p mysql/micro-blog/{data,backups}'
scp docker-stack.yml youruser@yourwebserver:~/
# Note: --with-registry-auth flag needed to send credentials to private container registry
ssh yourwebserver 'docker stack deploy -c docker-stack.yml davidado --with-registry-auth --resolve-image always --prune'
ssh yourwebserver 'docker service ls'

# Remove YAML file and unused images
ssh yourwebserver 'docker system prune -af && rm docker-stack.yml'