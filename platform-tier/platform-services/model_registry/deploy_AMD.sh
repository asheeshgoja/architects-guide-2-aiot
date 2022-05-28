
docker build -t model-registry-amd -f Dockerfile_AMD .
docker tag model-registry-amd:latest docker.35.238.247.144.nip.io:5000/model-registry-amd:latest
docker push docker.35.238.247.144.nip.io:5000/model-registry-amd:latest

# docker buildx build -t asheeshgoja/model-registry:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/model-registry:latest --platform linux/arm --push  .