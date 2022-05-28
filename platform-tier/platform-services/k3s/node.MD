**install**
On VM open firewall ports 
6443 and 5000

Make sure hostname resolves to IP address
Add entry in /etc/hosts

```bash
curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666 --tls-san 35.236.22.237 --node-external-ip=35.236.22.237" sh -


# curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--write-kubeconfig ~/.kube/config --write-kubeconfig-mode 666 --tls-san 35.236.22.237 --node-external-ip=35.236.22.237 --resolv-conf=/home/agoja/resolv.conf" sh -
```

<!-- curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable traefik --tls-san 35.236.22.237 --node-external-ip 35.236.22.237 --write-kubeconfig-mode 644" sh -s - -->


crictl info


**Uninstall**
```bash
/usr/local/bin/k3s-uninstall.sh
```

