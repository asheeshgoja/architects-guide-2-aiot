FROM nvcr.io/nvidia/l4t-tensorflow:r32.6.1-tf2.5-py3


RUN mkdir /tensorflow
WORKDIR /tensorflow

COPY train.py .
COPY loop.sh .
COPY training_dataset_1 .
COPY requirements.txt .

RUN pip3 install --upgrade setuptools
RUN python3 -m pip install --upgrade pip

# RUN pip3 install -U numpy
# RUN pip3 install -U scikit-learn

RUN pip3 install -r requirements.txt
ENV OPENBLAS_CORETYPE ARMV8

# CMD ["/tensorflow/loop.sh"]
CMD [ "python3", "/tensorflow/train.py" ]
