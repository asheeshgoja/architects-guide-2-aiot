import tensorflow as tf
import os
import requests
import datetime
import urllib  
import json
import ssl

from zipfile import ZipFile
from kafka import KafkaProducer


dir_path = os.path.dirname(os.path.realpath(__file__))

MODEL_DOWNLOAD_REGISTRY_URL = os.environ.get('MODEL_DOWNLOAD_REGISTRY_URL', 'https://localhost:8080/full')
MODEL_UPLOAD_REGISTRY_URL = os.environ.get('MODEL_UPLOAD_REGISTRY_URL', 'https://localhost:8080/uploadQuantizedModel')
KAFKA_BROKER = os.environ.get('KAFKA-BROKER', '35.236.22.237:32199')
CONTROL_TOPIC = os.environ.get('CONTROL-TOPIC', 'control-message')



def publishControlMessage(quantizedFileName):

    try:    
        fileName = os.path.basename(quantizedFileName)
        producer = KafkaProducer(bootstrap_servers=KAFKA_BROKER)
        json_data = {"command": "download-model", "payload" : fileName}
        message = json.dumps(json_data)
        bytesMessage = message.encode()
        producer.send(CONTROL_TOPIC, bytesMessage )        
    except Exception as x:
        print("failed to publish control message " + quantizedFileName)
        print(x)


def uploadModel(file_name):
    try:
        model_file = open(file_name, "rb")

        upload_url = MODEL_UPLOAD_REGISTRY_URL 

        test_response = requests.post(upload_url, files = {"file": model_file}, verify=False)

        if test_response.ok:
            print("Model Uploaded completed successfully to url " + upload_url)
            print(test_response.text)
        else:
            print("Model upload failed to url " + upload_url)

        publishControlMessage(file_name)
    except Exception as x:
        print("Model upload failed to url " + upload_url)
        print(x)


def quantizeModel(fileName):

    try:
        quantized_file_name = "{}/{}-model.tflite".format(dir_path,datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S') )

        converter = tf.lite.TFLiteConverter.from_saved_model(fileName)
        tflite_model = converter.convert()
        open(quantized_file_name, "wb").write(tflite_model)

        print("quantized downloaded model file : " + quantized_file_name )
        uploadModel(quantized_file_name)
    except Exception as x:
        print(x)



def unZipModelFiles(file_name):
    # opening the zip file in READ mode
    with ZipFile(file_name, 'r') as zip:
        # printing all the contents of the zip file
        zip.printdir()
        # for zip_info in zip.infolist():
        #     if zip_info.filename[-1] == '/':
        #         continue
        #     zip_info.filename = os.path.basename(zip_info.filename)
        #     zip.extract(zip_info, my_dir)


        model_dir = os.path.dirname(zip.infolist()[0].filename)

        # extracting all the files
        print('Extracting all the model files now...')

        zip.extractall(path=dir_path)
        return  dir_path + "/" + model_dir

def downloadAndQuantizeModelFiles():
    model_registry_url = MODEL_DOWNLOAD_REGISTRY_URL

    # use this for valid certs
    # context = ssl.create_default_context() 
    # context.load_cert_chain('/keys/ssh-publickey', '/keys/ssh-privatekey')
    # context.verify_mode = ssl.CERT_OPTIONAL

    # use this for self signed certs certs
    context = ssl._create_unverified_context() 

    opener = urllib.request.build_opener(urllib.request.HTTPSHandler(context=context))
    urllib.request.install_opener(opener)


    print("downloading model file from url : " + model_registry_url )
    
    urllib.request.urlretrieve(model_registry_url, dir_path + '/model_store_dir.xml' )

    from xml.dom import minidom
    xmldoc = minidom.parse(dir_path + '/model_store_dir.xml')
    itemlist = xmldoc.getElementsByTagName('a')

    for item in itemlist:
        latest_model_file_name = item.attributes['href'].value
        model_registry_url = MODEL_DOWNLOAD_REGISTRY_URL + '/' + latest_model_file_name
        downloaded_file_name = dir_path + '/' + latest_model_file_name
        urllib.request.urlretrieve(model_registry_url, downloaded_file_name)
        print("downloaded model file : " + downloaded_file_name  )
        model_dir = unZipModelFiles(downloaded_file_name)
        quantizeModel(model_dir)
        # quantizeModel("/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/quantization_module/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_module/2021-10-07-04:51:04-savedDir")




if __name__ == '__main__':
    downloadAndQuantizeModelFiles()