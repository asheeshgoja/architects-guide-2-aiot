
docker build -t model-registry .
docker tag model-registry:latest docker.35.238.247.144.nip.io:5000/model-registry:latest
docker push docker.35.238.247.144.nip.io:5000/model-registry:latest

# docker buildx build -t asheeshgoja/model-registry:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/model-registry:latest --platform linux/arm --push  .