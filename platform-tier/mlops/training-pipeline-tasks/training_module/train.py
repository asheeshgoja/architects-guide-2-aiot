# from turtle import mode
import ssl
import sklearn
# Important to keep sklearn as the first import for nvidia jetson nano and export env OPENBLAS_CORETYPE = ARMV8
# OPENBLAS_CORETYPE=ARMV8 python

import pandas
import tensorflow as tf
import os
import requests
import datetime
import urllib
from numpy import loadtxt

from zipfile import ZipFile

from tensorflow.python.keras import Sequential
from tensorflow.python.keras.layers import Dense
from tensorflow.python.keras.wrappers.scikit_learn import KerasClassifier
from tensorflow.python.keras import models


from sklearn.model_selection import train_test_split

dir_path = os.path.dirname(os.path.realpath(__file__))

MODEL_REGISTRY_URL = os.environ.get('MODEL_REGISTRY_URL', 'https://localhost:8080/uploadModel')
TRAINING_DATA_URL = os.environ.get('TRAINING_DATA_URL', 'https://localhost:8080/normalized_training_data/')
EPOCS = os.environ.get('EPOCS', '2')
BATCH_SIZE = os.environ.get('BATCH_SIZE', '32')


def uploadModel(file_name):

    try:
        # test_file = open("/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_module/model.h5", "rb")
        model_file = open(file_name, "rb")

        # test_url = "http://model-registry-service.architectsguide2aiot.svc.cluster.local:30007/uploadmodel"
        upload_url = MODEL_REGISTRY_URL  # "http://10.0.0.31:30007/uploadModel"

        test_response = requests.post(upload_url, files={"file": model_file}, verify=False)

        if test_response.ok:
            print("Model Uploaded completed successfully to url " + upload_url)
            print(test_response.text)
        else:
            print("Model upload failed to url " + upload_url)
            print(test_response)
    except Exception as x:
        print("Model upload failed to url " + upload_url)
        print(x)


def get_all_file_paths(directory):

    # initializing empty file paths list
    file_paths = []

    # crawling through directory and subdirectories
    for root, directories, files in os.walk(directory):
        for filename in files:
            # join the two strings in order to form the full filepath.
            filepath = os.path.join(root, filename)
            file_paths.append(filepath)

    # returning all file paths
    return file_paths


def zipSavedModelFiles(directory, zipFileName):
    # path to folder which needs to be zipped
    # directory = './python_files'

    # calling function to get all file paths in the directory
    file_paths = get_all_file_paths(directory)

    # printing the list of all files to be zipped
    print('Following model files will be zipped:')
    for file_name in file_paths:
        print(file_name)

    # writing files to a zipfile
    with ZipFile(zipFileName, 'w') as zip:
        # writing each file one by one
        for file in file_paths:
            zip.write(file)

    print('All tf model files zipped successfully!')


def trainModel_OLD():

    # dataframe = pandas.read_csv("/home/agoja/kubecon-2021-aiot-demo/apps/logistic_regression_model/training_module/training_dataset_1", header=None,skiprows=1)
    dataframe = pandas.read_csv("/home/agoja/kubecon-2021-aiot-ops-ref-app/apps/logistic_regression_model/training_module/training_dataset_1", header=None, skiprows=1)
    del dataframe[0]
    del dataframe[1]

    dataset = dataframe.values
    X = dataset[:, 0:4].astype(float)
    y = dataset[:, 4]

    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.33)

    model = Sequential()
    model.add(Dense(60, input_dim=4, activation='relu'))
    model.add(Dense(30, activation='relu'))
    model.add(Dense(10, activation='relu'))
    model.add(Dense(1, activation='sigmoid'))
    model.compile(loss='binary_crossentropy',
                  optimizer='adam', metrics=['accuracy'])
    model.fit(X_train, y_train, epochs=int(EPOCS),
              batch_size=int(BATCH_SIZE), verbose=0)

    # file_name = datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S')
    zip_file_name = "{}/{}-model.zip".format(
        dir_path, datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S'))
    save_dir_name = "{}/{}-savedDir".format(
        dir_path, datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S'))
    model.save(save_dir_name)
    zipSavedModelFiles(save_dir_name, zip_file_name)

    # converter = tf.lite.TFLiteConverter.from_keras_model(model)
    # tflite_model = converter.convert()
    # open(dir_path + "/" + file_name + "-model.tflite", "wb").write(tflite_model)

    loss, acc = model.evaluate(X_test, y_test, verbose=0)
    print('Test Accuracy: %.3f' % acc)

    uploadModel(zip_file_name)


def trainModel(trainingDataFile):

    dataset = loadtxt(trainingDataFile, delimiter=',', skiprows=1, usecols=range(2,7))
    X = dataset[:,0:4]
    y = dataset[:,4]

    model = Sequential()
    # This can be furhter optimized by hyperparameter tuning
    model.add(Dense(50, input_dim=4, activation='relu'))
    model.add(Dense(40, activation='relu'))
    model.add(Dense(1, activation='sigmoid'))
    model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
    model.fit(X, y, epochs=int(EPOCS),batch_size=int(BATCH_SIZE), verbose=0)

    # evaluate the model
    loss, acc = model.evaluate(X, y, verbose=0)
    print('Test Accuracy: %.3f' % acc)

#   freeze, zip and save the model
    zip_file_name = "{}/{}-model.zip".format(dir_path, datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S'))
    save_dir_name = "{}/{}-savedDir".format(dir_path, datetime.datetime.now().strftime('%Y-%m-%d-%H:%M:%S'))
    model.save(save_dir_name)
    zipSavedModelFiles(save_dir_name, zip_file_name)
    return zip_file_name


def downloadNormalizedTrainingData():
    # use this for valid certs
    # context = ssl.create_default_context() 
    # context.load_cert_chain('/keys/ssh-publickey', '/keys/ssh-privatekey')
    # context.verify_mode = ssl.CERT_OPTIONAL

    # use this for self signed certs certs
    context = ssl._create_unverified_context() 

    opener = urllib.request.build_opener(urllib.request.HTTPSHandler(context=context))
    urllib.request.install_opener(opener)


    training_data_url = TRAINING_DATA_URL

    print("downloading training data file from url : " + training_data_url)

    urllib.request.urlretrieve(
        training_data_url, dir_path + '/training_data_index.xml')

    from xml.dom import minidom
    xmldoc = minidom.parse(dir_path + '/training_data_index.xml')
    itemlist = xmldoc.getElementsByTagName('a')
    training_data_sets = []

    for item in itemlist:
        latest_model_file_name = item.childNodes[0].nodeValue
        training_data_url = TRAINING_DATA_URL +  latest_model_file_name
        downloaded_file_name = dir_path + '/' + latest_model_file_name
        urllib.request.urlretrieve(training_data_url, downloaded_file_name)
        training_data_sets.append(downloaded_file_name)
        print("downloaded training data file : " + downloaded_file_name)

    return training_data_sets


if __name__ == '__main__':
    training_data_sets = downloadNormalizedTrainingData()
    # modelFileName = trainModel_OLD(training_data_sets[0])
    # trainModel_OLD()
    # modelFileName = trainModel("/home/agoja/kubecon-2021-aiot-ops-ref-app/apps/logistic_regression_model/training_module/training_dataset_1")
    modelFileName = trainModel(training_data_sets[0])
    # modelFileName = trainModel("/home/agoja/kubecon-2021-aiot-ops-ref-app/apps/logistic_regression_model/training_module/normalized_training_data_2022-04-19T17:49:54:1650390594561.csv")
    uploadModel(modelFileName)
