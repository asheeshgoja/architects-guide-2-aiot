docker build -t extract_module .
docker tag extract_module:latest docker.35.238.247.144.nip.io:5000/extract_module:latest
docker push docker.35.238.247.144.nip.io:5000/extract_module:latest

# docker buildx build -t asheeshgoja/extract_module:latest --platform linux/arm64 --push  .
