#!/bin/bash

ssh -A deploy@172.25.30.46 '/bin/bash -s' << EOF
docker pull bahramkb/imageservice:latest
docker rmi bahramkb/imageservice:current || true
docker tag bahramkb/imageservice:latest bahramkb/imageservice:current
mkdir ~/swarm || true
cd ~/swarm
echo 'yes\n' | git clone git@github.com:bahramkarimi/imageservice.git ||
cd imageservice
git pull || true
docker stack deploy -c docker-compose.yml imageservice
EOF