FROM debian:buster

# FROM python:3.-7slim-bullseye


RUN mkdir /coral
WORKDIR /coral

RUN apt update && \
apt-get install curl gnupg ca-certificates -y

RUN echo "deb https://packages.cloud.google.com/apt coral-edgetpu-stable main" | tee /etc/apt/sources.list.d/coral-edgetpu.list

RUN echo "deb https://packages.cloud.google.com/apt coral-cloud-stable main" | tee /etc/apt/sources.list.d/coral-cloud.list

RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

RUN apt-get update


RUN apt-get update && \
    apt-get install python3  python3-pip -y                

RUN apt-get install python3-pycoral -y


COPY loop.sh ./

COPY infer_tflite_socket.py ./
# COPY model.tflite ./
COPY requirements.txt ./

# RUN pip3 install -r requirements.txt

EXPOSE 9898

CMD ["/coral/loop.sh"]
# CMD [ "python3", "/coral/infer_tflite_socket.py" ]