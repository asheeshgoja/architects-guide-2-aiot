# manually get the buildx image sha for asheeshgoja/edge-tpu-inference-engine:latest and update kubernetes_share_volumes.yaml

kubectl label nodes follower-node-coral-tpu1 tpuAccelerator=true
kubectl label nodes follower-node-coral-tpu2 tpuAccelerator=true
kubectl label nodes follower-node-coral-tpu3 tpuAccelerator=true

kubectl delete -f kubernetes_share_volumes.yaml -n architectsguide2aiot

sleep 5

kubectl apply -f kubernetes_share_volumes.yaml -n architectsguide2aiot
kubectl get pods -o wide -n architectsguide2aiot -w


# test
#  kubectl exec --stdin --tty -n architectsguide2aiot coral-python-deployment-  -- tail /coral/infer_tflite_socket.txt -f
#  kubectl exec --stdin --tty -n argo hello-world-g46kf -c main -- tail /coral/infer_tflite_socket.txt -f
# cd kubecon-2021-aiot-demo/apps/kafka_producer_client/go_console_producer/
# go run .

# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 36.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 39.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 45.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 52.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 59.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "rotation": 628.4, "temperature": 79.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:09:04", "rotation": 1634.7, "temperature": 145.4, "vibration":1235.9, "sound": 150.6}

# kubectl exec --stdin --tty coral-python-deployment-84ff577c55-qxz7c -n kafka  -- tail /coral/infer_tflite_socket.txt -f