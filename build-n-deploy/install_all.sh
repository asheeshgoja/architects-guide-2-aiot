# kubectl create ns architectsguide2aiot
# sleep 5

# # platform - services
# strimizi
../platform-tier/platform-services/strimizi/untaint.sh
../platform-tier/platform-services/strimizi/taint.sh
kubectl apply -f ../platform-tier/platform-services/strimizi/strimzi_kafka_operator.yaml  -n architectsguide2aiot
sleep 5
kubectl apply -f ../platform-tier/platform-services/strimizi/kafka-epheremal-single.yaml  -n architectsguide2aiot
echo "waiting for kafka cluster to start..."
kubectl wait kafka/architectsguide2aiot-cluster --for=condition=Ready --timeout=300s -n architectsguide2aiot
../platform-tier/platform-services/strimizi/untaint.sh

# # argo
kubectl apply -f  ../platform-tier/platform-services/argo_workflows/argo-deployment.yaml -n architectsguide2aiot
sleep 10
../platform-tier/platform-services/argo_workflows/patch.sh

# platform microservices
kubectl apply -f ../platform-tier/platform-services/mqtt_broker/kubernetes.yaml  -n architectsguide2aiot
kubectl wait --for=condition=ready pod -n architectsguide2aiot -l app=mqtt-broker

kubectl apply -f ../platform-tier/platform-services/model_registry/kubernetes.yaml  -n architectsguide2aiot
kubectl wait --for=condition=ready pod -n architectsguide2aiot -l app=model-registry

kubectl apply -f ../platform-tier/platform-services/training_datastore_μservice/kubernetes.yaml  -n architectsguide2aiot
kubectl wait --for=condition=ready pod -n architectsguide2aiot -l app=training-datastore

kubectl apply -f ../platform-tier/platform-services/mqtt-kafka-protocol-bridge/kubernetes.yaml  -n architectsguide2aiot
kubectl apply -f ../platform-tier/mlops/ingest_μservice/kubernetes.yaml  -n architectsguide2aiot
kubectl apply -f ../platform-tier/platform-services/device_registry/kubernetes.yaml  -n architectsguide2aiot

# # # workflow training pipeline DAG
# argo submit -n architectsguide2aiot --serviceaccount argo --watch ../platform-tier/mlops/argo-demo-dag/demo_dag.yaml

# # # edge tup inference modules
# kubectl apply -f  ../inference-tier/kubernetes.yaml -n architectsguide2aiot


# # kubectl apply -f  ../platform-tier/mlops/training-pipeline-tasks/quantization_module/kubernetes.yaml -n architectsguide2aiot
# # kubectl apply -f  ../platform-tier/mlops/training-pipeline-tasks/training_module/kubernetes.yaml -n architectsguide2aiot
# # kubectl apply -f  ../platform-tier/mlops/training-pipeline-tasks/extract_module/kubernetes.yaml -n architectsguide2aiot
# # kubectl apply -f  ../platform-tier/mlops/training-pipeline-tasks/validation_module/kubernetes.yaml -n architectsguide2aiot


# delete all data
# curl https://10.0.0.30:30008/deleteAllData
# curl https://10.0.0.30:30007/deleteAllData

# quay.io/strimzi/operator:0.25.0
# quay.io/strimzi/operator:0.27.0-rc1-0-arm64

# quay.io/strimzi/kafka:0.25.0-kafka-2.8.0
# quay.io/strimzi/kafka:0.27.1-0-kafka-2.8.0-arm64



# quay.io/strimzi/kafka:0.28.0-kafka-3.1.0
# quay.io/strimzi/kafka:0.28.0-kafka-3.1.0-arm64

# quay.io/strimzi/operator:0.28.0
# quay.io/strimzi/operator:0.28.0-arm64