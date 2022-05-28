get token from controller node
cat /var/lib/rancher/k3s/server/token


curl -sfL https://get.k3s.io | K3S_URL=https://35.236.22.237:6443 K3S_TOKEN=K1044479f717480f437d616a5b03f1940eadfbb017bee15af92572ed6932494ed44::server:7a08e3480e0545168e2b1b3a528ab80b sh -
<!-- curl -sfL https://get.k3s.io | K3S_URL=https://35.236.22.237:6443 K3S_TOKEN=K10b0226ae3db52c7613c5bd3861fef8981c110424cb0ab19532c70c30c7f7759a7::server:1d49b5a7e050b5c630f5de5ee2e8cdd5 K3S_RESOLV_CONF=/etc/resolv.conf sh - -->
	

status
crictl info


Agent uninstall
/usr/local/bin/k3s-agent-uninstall.sh

**Raspi 4**
sudo apt-get install vim
add cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1
to /boot/cmdline.txt
reboot


**nvidia jetson - remove docker**
```bash
dpkg -l | grep -i docker
sudo apt-get purge -y docker-engine docker docker.io docker-ce docker-ce-cli
sudo apt-get autoremove -y --purge docker-engine docker docker.io docker-ce  
sudo rm -rf /var/lib/docker /etc/docker
sudo rm /etc/apparmor.d/docker
sudo groupdel docker
sudo rm -rf /var/run/docker.sock
sudo rm -rf ~/.docker


Check
sudo iptables -L

reboot now
```


**label and taint nodes**
kubectl label nodes agentnode-coral-tpu1 tpuAccelerator=true
kubectl label nodes agentnode-coral-tpu2 tpuAccelerator=true
kubectl label nodes agentnode-coral-tpu3 tpuAccelerator=true
kubectl label nodes agentnode-nvidia-jetson gpuAccelerator=true



kubectl label nodes kubecon-aiot-control-node dedicated=Kafka

kubectl taint nodes agentnode-raspi1 dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-raspi2 dedicated=Kafka:NoSchedule
<!-- kubectl taint nodes agentnode-nvidia-jetson dedicated=Kafka:NoSchedule -->
kubectl taint nodes agentnode-coral-tpu1 dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-coral-tpu2 dedicated=Kafka:NoSchedule
kubectl taint nodes agentnode-coral-tpu3 dedicated=Kafka:NoSchedule




http://74.125.141.128/architectsguide2aiot-aiot-mlops-demo/agglomeration-tower1-cframe-shaded-pole_solvent_motor.csv
http://storage.googleapis.com/architectsguide2aiot-aiot-mlops-demo/agglomeration-tower1-cframe-shaded-pole_solvent_motor.csv

