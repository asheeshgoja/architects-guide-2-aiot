# docker rmi -f $(docker images | grep 'go_test_consumer')
# docker build -t go_test_consumer -f Dockerfile_consumer .
# docker buildx build -t go_test_consumer -f Dockerfile_consumer . --platform linux/arm

# docker buildx build -t asheeshgoja/go_test_consumer:latest --platform linux/arm64 --push  -f Dockerfile_consumer .

# docker tag go_test_consumer:latest docker.104.197.50.43.nip.io:5000/go_test_consumer:latest
# docker push docker.35.238.247.144.nip.io:5000/go_test_consumer_small:latest

#######


docker build -t go_test_consumer_small -f Dockerfile .
docker tag go_test_consumer_small:latest docker.35.238.247.144.nip.io:5000/go_test_consumer_small:latest
docker push docker.35.238.247.144.nip.io:5000/go_test_consumer_small:latest