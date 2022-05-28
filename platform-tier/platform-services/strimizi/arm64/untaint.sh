# kubectl taint nodes agentnode-raspi1 dedicated=Kafka:NoSchedule-
# kubectl taint nodes agentnode-nvidia-jetson dedicated=Kafka:NoSchedule-
kubectl taint nodes agentnode-coral-tpu1 dedicated=Kafka:NoSchedule-
kubectl taint nodes agentnode-coral-tpu2 dedicated=Kafka:NoSchedule-

kubectl taint nodes kubecon-aiot-control-node dedicated=Kafka:NoSchedule-
