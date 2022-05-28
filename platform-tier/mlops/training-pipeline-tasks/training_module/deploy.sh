
docker build -t training-module .
docker tag training-module:latest docker.35.238.247.144.nip.io:5000/training-module:latest
docker push docker.35.238.247.144.nip.io:5000/training-module:latest

# docker build -t training-module .
# docker tag training-module:latest asheeshgoja/training-module:latest
# docker push asheeshgoja/training-module:latest

# docker buildx build -t asheeshgoja/training-module:latest --platform linux/arm64 --push  .
