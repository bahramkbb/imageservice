#!/bin/bash

echo 'Please enter server IP address:'
read MASTER_IP

ssh -A deploy@$MASTER_IP /bin/bash << EOF
docker pull bahramkb/imageservice:latest
mkdir ~/swarm || true
cd ~/swarm
echo 'yes\n' | git clone git@github.com:bahramkarimi/imageservice.git || true
cd imageservice
git pull || true
docker stack deploy -c docker-compose.yml imageservice
EOF