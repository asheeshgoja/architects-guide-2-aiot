import requests

# test_file = open("/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_module/model.h5", "rb")
model_file = open("/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_module/model.tflite", "rb")

test_url = "http://localhost:8080/uploadmodel"

test_response = requests.post(test_url, files = {"file": model_file})
# test_response = requests.post(test_url, files = {"name": "model.tflite", "file": model_file})


if test_response.ok:
    print("Upload completed successfully!")
    print(test_response.text)
else:
    print("Something went wrong!")
