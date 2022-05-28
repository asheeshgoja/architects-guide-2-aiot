# kubectl delete -f quick-start-postgres.yaml -n argo

# sleep 10

# kubectl apply -n argo -f quick-start-postgres.yaml

kubectl create ns argo
kubectl delete -n architectsguide2aiot -f install-v3-1-1.yaml

kubectl apply -n architectsguide2aiot -f install-v3-1-1.yaml



kubectl patch configmap/workflow-controller-configmap \
-n architectsguide2aiot \
--type merge \
-p '{"data":{"containerRuntimeExecutor":"k8sapi"}}'


# Now using nodeport @ port  30042
# kubectl -n architectsguide2aiot port-forward svc/argo-server 2746:2746
# forward port 2746 in VS Code 
# https://127.0.0.1:2746
# https://35.236.22.237:30042/
# https://10.0.0.30:30042/
# sudo lsof -i -P -n | grep LISTEN

##UI Aauth
kubectl -n architectsguide2aiot exec argo-server- -- argo auth token


argo submit -n architectsguide2aiot --serviceaccount argo --watch demo_DAG.yaml

argo submit -n architectsguide2aiot --serviceaccount argo --watch steps.yaml 

# kubectl exec --stdin --tty -n architectsguide2aiot kubecon-aiotdemo-dag- -c edge-tpu-inference-engine -- tail /coral/infer_tflite_socket.txt -f












# //ARGO EVENTS

kubectl create ns argo-events
# Install Argo Events
# kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install.yaml
kubectl apply -f install_argo_events.yaml
# Deploy the Event Bus

# $ kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/eventbus/native.yaml
kubectl apply -n argo-events -f native.yaml


# kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/master/examples/rbac/sensor-rbac.yaml
kubectl apply -n argo-events -f sensor-rbac.yaml


# kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/event-sources/webhook.yaml
kubectl apply -n argo-events -f event-sources-webhook.yaml

# kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/sensors/webhook.yaml
kubectl apply -n argo-events -f sensors-webhook.yaml 

kubectl -n argo-events port-forward webhook-eventsource-ktmd6-7dc686497d-vvps2 12000:12000

http://localhost:12000/example


# //
kubectl delete -n argo-events -f sensors-webhook.yaml 
kubectl delete -n argo-events -f event-sources-webhook.yaml
kubectl delete -n argo-events -f sensor-rbac.yaml
kubectl delete -n argo-events -f native.yaml
kubectl delete -f install_argo_events.yaml
