
docker build -t protocol_bridge .
docker tag protocol_bridge:latest docker.35.238.247.144.nip.io:5000/protocol_bridge:latest
docker push docker.35.238.247.144.nip.io:5000/protocol_bridge:latest

# docker buildx build -t asheeshgoja/protocol_bridge:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/protocol_bridge:latest --platform linux/arm --push  .