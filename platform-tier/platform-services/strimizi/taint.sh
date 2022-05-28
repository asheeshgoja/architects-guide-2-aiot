kubectl label nodes kubecon-aiot-control-node dedicated=Kafka
kubectl label nodes agentnode-coral-tpu1 tpuAccelerator=true
kubectl label nodes agentnode-coral-tpu2 tpuAccelerator=true
# kubectl label nodes agentnode-coral-tpu3 tpuAccelerator=true
kubectl label nodes agentnode-nvidia-jetson gpuAccelerator=true

kubectl taint nodes agentnode-raspi1 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-raspi2 dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-nvidia-jetson dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-coral-tpu1 dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-coral-tpu2 dedicated=Kafka:NoSchedule
# kubectl taint nodes agentnode-coral-tpu3 dedicated=Kafka:NoSchedule

# kubectl label nodes agentnode-raspi1 controlnode=active