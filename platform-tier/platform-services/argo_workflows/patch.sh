kubectl patch configmap/workflow-controller-configmap \
-n architectsguide2aiot \
--type merge \
-p '{"data":{"containerRuntimeExecutor":"k8sapi"}}'
