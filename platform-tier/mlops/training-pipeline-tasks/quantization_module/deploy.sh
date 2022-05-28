
docker build -t quantize-module .
docker tag quantize-module:latest docker.35.238.247.144.nip.io:5000/quantize-module:latest
docker push docker.35.238.247.144.nip.io:5000/quantize-module:latest

# docker build -t quantize-module .
# docker tag quantize-module:latest asheeshgoja/quantize-module:latest
# docker push asheeshgoja/quantize-module:latest

# docker buildx build -t asheeshgoja/quantize-module:latest --platform linux/arm64 --push  .
