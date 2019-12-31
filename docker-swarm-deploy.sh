#!/bin/bash
#

ssh jado 'mkdir -p mysql/micro-user/{data,backups};mkdir -p mysql/micro-blog/{data,backups}'
scp docker-stack.yml jason@jado:~/
# Note: --with-registry-auth flag needed to send credentials to private container registry
ssh jado 'docker stack deploy -c docker-stack.yml davidado --with-registry-auth --resolve-image always --prune'
ssh jado 'docker service ls'

# Remove YAML file and unused images
# ssh jado 'docker system prune -af && rm docker-stack.yml'