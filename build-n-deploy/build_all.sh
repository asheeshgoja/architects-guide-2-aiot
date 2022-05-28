# clear
# docker rm -vf $(docker ps -a -q)
# docker rmi -f $(docker images -a -q)
# sleep 5
# docker run -d -p 5000:5000 --restart=always --name registry registry:2
# sleep 5

cd ../platform-tier/platform-services/model_registry
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/platform-services/device_registry
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/platform-services/training_datastore_μservice
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/mlops/ingest_μservice
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/platform-services/mqtt-kafka-protocol-bridge
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/platform-services/mqtt-broker
./deploy.sh
sleep 5
cd ../../../build-n-deploy/

cd ../platform-tier/mlops/training-pipeline-tasks/extract_module
./deploy.sh
sleep 5
cd ../../../../build-n-deploy/

cd ../platform-tier/mlops/training-pipeline-tasks/quantization_module
./deploy.sh
sleep 
cd ../../../../build-n-deploy/

cd ../platform-tier/mlops/training-pipeline-tasks/training_module
./deploy.sh
sleep 5
cd ../../../../build-n-deploy/

cd ../platform-tier/mlops/training-pipeline-tasks/validation_module
./deploy.sh
sleep 5
cd ../../../../build-n-deploy/

cd ../inference-tier
./build.sh
sleep 5


