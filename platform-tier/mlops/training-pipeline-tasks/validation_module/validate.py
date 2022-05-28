import ssl
import tensorflow_data_validation as tfdv
import sys
import os
import urllib
import json 
from kafka import KafkaProducer


TRAINING_DATA_URL = os.environ.get('TRAINING_DATA_URL', 'https://localhost:8080/normalized_training_data/')
# TRAINING_DATA_URL = os.environ.get('TRAINING_DATA_URL', 'http://73.252.176.163:30007/normalized_training_data/')
dir_path = os.path.dirname(os.path.realpath(__file__))
KAFKA_BROKER = os.environ.get('KAFKA-BROKER', '35.236.22.237:32199')
CONTROL_TOPIC = os.environ.get('CONTROL-TOPIC', 'control-message')
DRIFT_THRESHOLD = os.environ.get('DRIFT_THRESHOLD', 6)
 


def validate(training_data_sets):

    for training_data_set in  training_data_sets:

        print("generating  statistics from csv file {} ".format(training_data_set))
        train_stats = tfdv.generate_statistics_from_csv(data_location=training_data_set)
       
        print("infering  schema from csv file {} ".format(training_data_set))
        schema = tfdv.infer_schema(train_stats)

        print("validating statistics for csv file {} ".format(training_data_set))
        anomalies = tfdv.validate_statistics(statistics=train_stats, schema=schema)
        # data = tfdv.validate_statistics(statistics=train_stats, schema=schema).data
        # tfdv.display_anomalies(anomalies)
        print("printing anomalies for validation")
        print(anomalies)

        if len(anomalies.baseline.feature) > DRIFT_THRESHOLD :
            publishControlMessage(training_data_set)



def downloadTrainingData():
    training_data_url = TRAINING_DATA_URL

    # use this for valid certs
    # context = ssl.create_default_context() 
    # context.load_cert_chain('/keys/ssh-publickey', '/keys/ssh-privatekey')
    # context.verify_mode = ssl.CERT_OPTIONAL

    # use this for self signed certs certs
    context = ssl._create_unverified_context() 

    opener = urllib.request.build_opener(urllib.request.HTTPSHandler(context=context))
    urllib.request.install_opener(opener)


    print("downloading training data file from url : " + training_data_url )
    
    urllib.request.urlretrieve(training_data_url, dir_path + '/training_data_index.xml')

    from xml.dom import minidom
    xmldoc = minidom.parse(dir_path + '/training_data_index.xml')
    itemlist = xmldoc.getElementsByTagName('a')
    training_data_sets = []

    for item in itemlist:
        latest_model_file_name = item.attributes['href'].value
        training_data_url = TRAINING_DATA_URL + '/' + latest_model_file_name
        downloaded_file_name = dir_path + '/' + latest_model_file_name
        urllib.request.urlretrieve(training_data_url, downloaded_file_name)
        training_data_sets.append(downloaded_file_name)
        print("downloaded training data file : " + downloaded_file_name)

    return training_data_sets


def publishControlMessage(trainingFileName):

    try:    
        fileName = os.path.basename(trainingFileName)
        producer = KafkaProducer(bootstrap_servers=KAFKA_BROKER)
        json_data = {"command": "train-model", "payload" : fileName}
        message = json.dumps(json_data)
        bytesMessage = message.encode()
        producer.send(CONTROL_TOPIC, bytesMessage )        
    except Exception as x:
        print("failed to publish control message " + trainingFileName)
        print(x)



if __name__ == '__main__':
    print("cmd line args : {}".format(sys.argv) )
    training_data_sets = downloadTrainingData()
    validate(training_data_sets)