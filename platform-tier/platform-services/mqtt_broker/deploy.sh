
docker build -t go_amr64_mqtt_broker .
docker tag go_amr64_mqtt_broker:latest docker.35.238.247.144.nip.io:5000/go_amr64_mqtt_broker:latest
docker push docker.35.238.247.144.nip.io:5000/go_amr64_mqtt_broker:latest

# docker buildx build -t asheeshgoja/go_amr64_mqtt_broker:latest --platform linux/arm64 --push  .

# raspi
# docker buildx build -t asheeshgoja/go_amr64_mqtt_broker:latest --platform linux/arm --push  .