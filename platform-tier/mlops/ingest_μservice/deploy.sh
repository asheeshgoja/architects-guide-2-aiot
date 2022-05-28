docker build -t ingest_service .
docker tag ingest_service:latest docker.35.238.247.144.nip.io:5000/ingest_service:latest
docker push docker.35.238.247.144.nip.io:5000/ingest_service:latest

# docker buildx build -t asheeshgoja/ingest_service_golang:latest --platform linux/arm64 --push  .
