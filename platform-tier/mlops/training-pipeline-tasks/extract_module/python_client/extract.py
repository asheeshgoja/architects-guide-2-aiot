from google.cloud import storage
import os
import requests
import datetime
import urllib


dir_path = os.path.dirname(os.path.realpath(__file__))
TRAINING_DATA_UPLOAD_REGISTRY_URL = os.environ.get('TRAINING_DATA_UPLOAD_REGISTRY_URL', 'http://localhost:8080/uploadTrainingData')
GCP_BUCKET = os.environ.get("GCP_BUCKET","architectsguide2aiot-aiot-mlops-demo")

def uploadTrainingDataToModelRegistry(file_name):
    try:
        model_file = open(file_name, "rb")
        upload_url = TRAINING_DATA_UPLOAD_REGISTRY_URL 
        test_response = requests.post(upload_url, files = {"file": model_file})

        if test_response.ok:
            print("Model Uploaded completed successfully to url " + upload_url)
            print(test_response.text)
        else:
            print("Model upload failed to url " + upload_url)
    except Exception as x:
        print("Model upload failed to url " + upload_url)
        print(x)


# def upload_to_gcp_storage(gcp_bucket, source_file_name, destination_blob_name):
#     storage_client = storage.Client.from_service_account_json(dir_path + '/cloud_key.json')
#     bucket = storage_client.bucket(gcp_bucket)
#     blob = bucket.blob(destination_blob_name)
#     blob.upload_from_filename(source_file_name)

def download_from_gcp_storage(gcp_bucket):
    try:
        gcp_storage_file_name = "agglomeration-tower1-cframe-shaded-pole_solvent_motor.csv"
        downloaded_training_file = "{}/{}-dataset.csv".format(dir_path,datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S') )

        storage_client = storage.Client.from_service_account_json(dir_path + '/cloud_key.json')
        blob = storage_client.bucket(gcp_bucket).blob(gcp_storage_file_name)
        blob.download_to_filename(downloaded_training_file)
        print("Downloaded gcp cloud storage file {} for local folder at {}".format(gcp_storage_file_name, downloaded_training_file))
        return downloaded_training_file

    except Exception as x:
        print(x)


# def downlodURL():

#     model_registry_url =  "http://storage.googleapis.com/architectsguide2aiot-aiot-mlops-demo/agglomeration-tower1-cframe-shaded-pole_solvent_motor.csv"

#     print("downloading model file from url : " + model_registry_url )
#     downloaded_training_file = "{}/{}-dataset.csv".format(dir_path,datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S') )
    
#     urllib.request.urlretrieve(model_registry_url, downloaded_training_file)
    
#     with open(downloaded_training_file, 'r') as f:
#         print(f.read())
    



if __name__ == '__main__':
    downloaded_training_file = download_from_gcp_storage(GCP_BUCKET)
    uploadTrainingDataToModelRegistry(downloaded_training_file)
    # downlodURL()



