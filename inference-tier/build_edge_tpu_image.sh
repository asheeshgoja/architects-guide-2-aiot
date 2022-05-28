#docker build -t edge-tpu-inference-engine -f Dockerfile .
#docker tag edge-tpu-inference-engine:latest docker.35.238.247.144.nip.io:5000/edge-tpu-inference-engine:latest
#docker push docker.35.238.247.144.nip.io:5000/edge-tpu-inference-engine:latest

cd edge_tpu_tflite_inference_engine/
# cp ../logistic_regression_model/training_module/model.tflite .
docker buildx build -t asheeshgoja/edge-tpu-inference-engine:latest --platform linux/arm64 --push  .


docker pull docker.io/asheeshgoja/edge-tpu-inference-engine:latest
docker tag docker.io/asheeshgoja/edge-tpu-inference-engine:latest docker.35.238.247.144.nip.io:5000/edge-tpu-inference-engine:latest
docker push docker.35.238.247.144.nip.io:5000/edge-tpu-inference-engine:latest