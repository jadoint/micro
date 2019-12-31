#!/bin/bash
#

# Initialize Docker swarm
ssh yourwebserver 'docker swarm init'

# Install private files on all servers
scp -r private yourwebserver:~/