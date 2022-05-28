FROM nvcr.io/nvidia/l4t-tensorflow:r32.6.1-tf2.5-py3


RUN mkdir /tensorflow
WORKDIR /tensorflow

COPY quantize.py .
COPY loop.sh .
COPY requirements.txt .

# RUN pip3 install --upgrade setuptools
# RUN python3 -m pip install --upgrade pip

RUN pip3 install -r requirements.txt

# CMD ["/tensorflow/loop.sh"]
CMD [ "python3", "/tensorflow/quantize.py" ]
