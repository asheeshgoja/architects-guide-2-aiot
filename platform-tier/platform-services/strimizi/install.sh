#!/bin/bash

# echo architectsguide2aiot

# if [[ -n architectsguide2aiot ]]; then
    # sed 's/namespace: kafka/namespace: architectsguide2aiot/' 02_strimzi_kafka_operator.yml

# kubectl label nodes kubecon-aiot-control-node dedicated=Kafka


# kubectl taint nodes agentnode-raspi1 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-raspi2 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-nvidia-jetson dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-coral-tpu1 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-coral-tpu2 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-coral-tpu3 dedicated=Kafka:NoSchedule

./taint.sh


kubectl create ns architectsguide2aiot 

# kubectl delete -f kafka-epheremal-single.yaml -n architectsguide2aiot 
# kubectl delete -f 02_strimzi_kafka_operator.yml -n architectsguide2aiot 


./taint.sh

kubectl apply -f 02_strimzi_kafka_operator.yml -n architectsguide2aiot 

sleep 5

kubectl apply -f kafka-epheremal-single.yaml -n architectsguide2aiot 

./untaint.sh

echo "waiting for kafka cluster to start"
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n architectsguide2aiot 
# else
#     echo "argument error, please specify namespace" 
# fi
