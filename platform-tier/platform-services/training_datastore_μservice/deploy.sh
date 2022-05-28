
docker build -t training-datastore .
docker tag training-datastore:latest docker.35.238.247.144.nip.io:5000/training-datastore:latest
docker push docker.35.238.247.144.nip.io:5000/training-datastore:latest

# docker buildx build -t asheeshgoja/training-datastore:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/training-datastore:latest --platform linux/arm --push  .