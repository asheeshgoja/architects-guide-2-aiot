kubectl get namespace longhorn-system -o json >tmp.json

curl -k -H "Content-Type: application/json" -X PUT --data-binary @tmp.json http://127.0.0.1:8001/api/v1/namespaces/longhorn-system/finalize

# Replica Zone Level Soft Anti-Affinity

//kubectl patch -n architectsguide2aiot pv artifacts-registry-volmume   -p '{"metadata":{"finalizers":null}}'
//kubectl patch  pv artifacts-registry-volmume   -p '{"metadata":{"finalizers":null}}'


 vi /etc/systemd/system/k3s.service.env
K3S_RESOLV_CONF=/etc/resolv.conf
sudo systemctl status k3s

//
create volume failes on raspi
create volume with iSCSI option first
attach it to raspi
detach and delete
recreate with Block device
kubecon

//
kubectl get volumeattachment
kubectl get volumeattachments.storage.k8s.io csi-7
kubectl edit  volumeattachments.storage.k8s.io csi-
// remove finalizeer section



//
kubectl get  pv,pvc
kubectl get  pv,pvc -n architectsguide2aiot
kubectl delete  pv artifacts-registry-volmume
kubectl delete -n architectsguide2aiot pvc artifacts-registry-volmume
kubectl edit -n architectsguide2aiot pvc artifacts-registry-volmume-x
kubectl edit  pv artifacts-registry-volmume