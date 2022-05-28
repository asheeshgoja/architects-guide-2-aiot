
docker build -t device-registry .
docker tag device-registry:latest docker.35.238.247.144.nip.io:5000/device-registry:latest
docker push docker.35.238.247.144.nip.io:5000/device-registry:latest

# docker buildx build -t asheeshgoja/device-registry:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/device-registry:latest --platform linux/arm --push  .