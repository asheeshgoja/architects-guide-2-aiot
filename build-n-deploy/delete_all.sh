# platform - services
kubectl delete -f ../platform-tier/platform-services/strimizi/kafka-epheremal-single.yaml  -n architectsguide2aiot
kubectl delete -f ../platform-tier/platform-services/strimizi/strimzi_kafka_operator.yaml  -n architectsguide2aiot
kubectl delete -f  ../platform-tier/platform-services/argo_workflows/argo-deployment.yaml -n architectsguide2aiot

kubectl delete -f ../platform-tier/platform-services/mqtt-kafka-protocol-bridge/kubernetes.yaml  -n architectsguide2aiot
kubectl delete -f ../platform-tier/platform-services/training_datastore_μservice/kubernetes.yaml  -n architectsguide2aiot
kubectl delete -f ../platform-tier/mlops/ingest_μservice/kubernetes.yaml  -n architectsguide2aiot
kubectl delete -f ../platform-tier/platform-services/model_registry/kubernetes.yaml  -n architectsguide2aiot

kubectl delete -f ../platform-tier/platform-services/device_registry/kubernetes.yaml  -n architectsguide2aiot
kubectl delete -f ../platform-tier/platform-services/mqtt_broker/kubernetes.yaml  -n architectsguide2aiot

kubectl delete -f  ../inference-tier/kubernetes.yaml -n architectsguide2aiot


