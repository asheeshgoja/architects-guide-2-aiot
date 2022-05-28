
docker build -t validation_module .
docker tag validation_module:latest docker.35.238.247.144.nip.io:5000/validation_module:latest
docker push docker.35.238.247.144.nip.io:5000/validation_module:latest

# docker build -t validation_module .
# docker tag validation_module:latest asheeshgoja/validation_module:latest
# docker push asheeshgoja/validation_module:latest

# docker buildx build -t asheeshgoja/validation_module:latest --platform linux/arm --push  .
