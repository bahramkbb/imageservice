#!/bin/bash
habitus
docker push bahramkb/imageservice

#ssh deploy@SWARMIP 
#docker-machine ssh cluster01 << EOF
#docker pull bahramkb/imageservice:latest
#docker rmi bahramkb/imageservice:current || true
#docker tag bahramkb/imageservice:latest bahramkb/imageservice:current
#EOF

docker-machine ssh cluster01 << EOF
mkdir /swarm || true
cd /swarm
git clone https://www.github.com/bahramkb/imageservice || true
git pull || true
docker stack deploy -c docker-compose.yml imageservice
EOF