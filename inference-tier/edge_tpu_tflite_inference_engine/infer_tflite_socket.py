import json
import socket
# import requests
import os
import ssl
import time
from types import SimpleNamespace
import numpy as np
import json
import urllib.request
from xml.dom import minidom


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'





#uncomment this section to test on a non edge tpu device
# import numpy as np
# import tensorflow as tf
# import numpy as np
# import json



#uncomment this section to test on a  edge tpu device
from pycoral.adapters import classify
from pycoral.adapters import common
from pycoral.utils.dataset import read_label_file
from pycoral.utils.edgetpu import make_interpreter
# from periphery import GPIO

# ledB = GPIO("/dev/gpiochip2", 9, "out")  # 16
# ledG = GPIO("/dev/gpiochip4", 10, "out") # 18
# ledR = GPIO("/dev/gpiochip4", 12, "out") # 22

def getInterpreter():
    # interpreter = tf.lite.Interpreter(model_path= dir_path + '/' + latest_model_file_name)
    interpreter = make_interpreter(dir_path + '/' + latest_model_file_name)
    
    return interpreter


dir_path = os.path.dirname(os.path.realpath(__file__))
SIDECAR_PORT = os.environ.get('SIDECAR_PORT', '9898')
MODEL_REGISTRY_URL = os.environ.get('MODEL_REGISTRY_URL', 'http://localhost:8080/quantized')


def logToFile(msg):
    
    print(msg )
    file1 = open("infer_tflite_socket.txt","a")
    file1.write(msg + "\n")
    file1.close()

def basic_consume_loop(latest_model_file_name):
    try:

        interpreter = getInterpreter()

        HOST = '127.0.0.1'  # Standard loopback interface address (localhost)
        PORT = int(SIDECAR_PORT)      # Port to listen on (non-privileged ports are > 1023)

        logToFile("Starting socket on {}:{}".format(HOST,PORT) )


        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind((HOST, PORT))
            s.listen()

            while True:
                try:
                    conn, addr = s.accept()
                    with conn:
                        logToFile('Connected by ' + addr[0] )
                        while True:
                            data = conn.recv(1024).decode()
                            logToFile("Received message from go client{}".format(data))
                            if not data:
                                break

                            ret = infer(interpreter, data)
                            strRet = "inference value from the edge tpu = {}".format(ret)
                            conn.sendall(strRet.encode())
                except Exception as inst:
                    print(inst)
            
    except Exception as inst:
        print(inst)


def shutdown():
    running = False


def getHostnameAndIpAddress():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    ipAdd = s.getsockname()[0]
    hostName = socket.getfqdn()
    return (ipAdd, hostName)


def infer(interpreter, msg):

       
    interpreter.allocate_tensors()

    # msg_fault_0 = """{"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 36.6, "vibration": 163.5, "sound": 50.0}"""
    # msg_fault_1 = """{"deviceID": "1", "timeStamp": "2021-09-14 23:09:04", "current": 1634.7, "temperature": 145.4, "vibration":1235.9, "sound": 150.6}"""

    # msg = msg_fault_1

    try:
        # current = np.float32(json.loads(msg)["current"])
        # temperature = np.float32(json.loads(msg)["temperature"])
        # vibration = np.float32(json.loads(msg)["vibration"])
        # sound = np.float32(json.loads(msg)["sound"])

        sensor_data = json.loads(msg, object_hook=lambda d: SimpleNamespace(**d))
        sensor_data_arr = [np.float32(sensor_data.current), np.float32(sensor_data.temperature), np.float32(sensor_data.vibration), np.float32(sensor_data.sound)]

        np_arr_64 = np.array(sensor_data_arr)
        np_arr_f32 = np_arr_64.astype(np.float32)

        inp_details = interpreter.get_input_details()
        out_details = interpreter.get_output_details()

        interpreter.set_tensor(inp_details[0]['index'], [sensor_data_arr])

        interpreter.invoke()

        output_details = interpreter.get_output_details()
        predictions = interpreter.get_tensor(output_details[0]['index'])

        output_index = interpreter.get_output_details()[0]["index"]
        ten = interpreter.get_tensor(output_index)

        inference_val = float(('%f' % ten))

        hn_ip_tuple = getHostnameAndIpAddress()

        logToFile(bcolors.BOLD +
              "Edge TPU logistic regression inference value : {}".format(inference_val))

        if 0.0 <= inference_val <= 0.5:
            blinkLed(ledG)
            logToFile(bcolors.ENDC + "No Fault detected on IPM-SynRM motor. Edge Tpu Node {1} , IP Address {0} ".format(
                hn_ip_tuple[0], hn_ip_tuple[1]))
        elif 0.5 <= inference_val <= 0.8:
            blinkLed(ledB)
            logToFile(bcolors.WARNING + "Potential conditions for an eminent IPM-SynRM motor failure. Edge Tpu Node {1}, IP Address {0} ".format(
                hn_ip_tuple[0], hn_ip_tuple[1]))
        elif 0.8 <= inference_val <= 1:
            blinkLed(ledR)
            logToFile(bcolors.FAIL + "Conditions for an immediate IPM-SynRM motor failure, activate safety protocols. Edge Tpu Node {1}, IP Address {0} ".format(
                hn_ip_tuple[0], hn_ip_tuple[1]))

        logToFile(bcolors.ENDC + "")

        return inference_val

    except Exception as inst:
        logToFile(inst.msg)
        return inst.msg

def blinkLed(led):
    # led.write(True)
    time.sleep(2)
    # led.write(False)

def downloadModelFile():
    model_registry_url = MODEL_REGISTRY_URL
    # use this for valid certs
    # context = ssl.create_default_context() 
    # context.load_cert_chain('/keys/ssh-publickey', '/keys/ssh-privatekey')
    # context.verify_mode = ssl.CERT_OPTIONAL

    # use this for self signed certs certs
    context = ssl._create_unverified_context() 

    opener = urllib.request.build_opener(urllib.request.HTTPSHandler(context=context))
    urllib.request.install_opener(opener)

    logToFile("downloading model file from url : " + model_registry_url )
    
    urllib.request.urlretrieve(model_registry_url, dir_path + '/model_store_dir.xml')


    # logToFile("downloading model file from url : " + model_registry_url )
    # r = requests.get(model_registry_url, allow_redirects=True)
    # open(dir_path + '/model_store_dir.xml', 'wb').write(r.content)
    # r.close()

    from xml.dom import minidom
    xmldoc = minidom.parse(dir_path + '/model_store_dir.xml')
    itemlist = xmldoc.getElementsByTagName('a')
    l = len(itemlist) #last in the list is the latest
    latest_model_file_name = itemlist[l-1].attributes['href'].value

    model_registry_url = MODEL_REGISTRY_URL + '/' + latest_model_file_name
    # r = requests.get(model_registry_url, allow_redirects=True)
    # open(dir_path + '/' + latest_model_file_name , 'wb').write(r.content)
    # r.close()
    urllib.request.urlretrieve(model_registry_url, dir_path + '/' + latest_model_file_name)

    logToFile("downloaded model file : " + latest_model_file_name  )

    return latest_model_file_name


if __name__ == '__main__':
    latest_model_file_name = downloadModelFile()
    basic_consume_loop(latest_model_file_name)


# kubectl run -n kafka kafka-producer -ti --image=strimzi/kafka:0.20.0-rc1-kafka-2.6.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --broker-list my-cluster-kafka-bootstrap.kafka:9092 --topic my-topic

# kubectl run -n kafka kafka-consumer -ti --image=strimzi/kafka:0.20.0-rc1-kafka-2.6.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap.kafka:9092 --topic my-topic --from-beginning


# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 36.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 39.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 45.6, "vibration": 163.5, "sound": 50.0}
# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 52.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 59.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 628.4, "temperature": 79.6, "vibration": 163.5, "sound": 50.0}

# {"deviceID": "1", "timeStamp": "2021-09-14 23:09:04", "current": 1634.7, "temperature": 145.4, "vibration":1235.9, "sound": 150.6}

# kubectl exec --stdin --tty coral-python-deployment- -n architectsguide2aiot  -- tail /coral/infer_tflite_socket.txt -f

# kubectl exec --stdin --tty kubecon-aiotdemo-dag-  -c edge-tpu-inference-engine  -n architectsguide2aiot  -- tail /coral/infer_tflite_socket.txt -f


# {"deviceID": "1", "timeStamp": "2021-09-14 23:07:59", "current": 6.4, "temperature": 36.6, "vibration": 163.5, "sound": 20971516.0}

# {"chipID": "14339204", "current": 628.29, "temperature": 36.95, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 628.29, "temperature": 65.95, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 628.29, "temperature": 68.95, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 628.29, "temperature": 70.95, "vibration": 160.00, "sound": 20.97}


# GREEN
# {"chipID": "14339204", "current": 228.55, "temperature": 21.65, "vibration": 167.00, "sound": 20.97}
# {"chipID": "14339204", "current": 627.93, "temperature": 85.65, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 627.93, "temperature": 95.65, "vibration": 160.00, "sound": 20.97}

# BLUE
# {"chipID": "14339204", "current": 627.93, "temperature": 110.65, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 627.93, "temperature": 115.5, "vibration": 160.00, "sound": 20.97}

# RED
# {"chipID": "14339204", "current": 627.93, "temperature": 160.5, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 627.93, "temperature": 160.5, "vibration": 160.00, "sound": 20.97}
# {"chipID": "14339204", "current": 627.93, "temperature": 195.65, "vibration": 160.00, "sound": 20.97}

#  {"chipID": "33", "current": "6.4","temperature": "36.6", "vibration": "163.5", "sound": "20971516.0"}
#  {"chipID": "3433", "current": "0.0","temperature": "145.4", "vibration": "1235.9", "sound": "150.6"}
