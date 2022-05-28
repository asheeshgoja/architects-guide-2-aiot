# # FROM balenalib/raspberry-pi-python:3-buster@sha256:68a082121f3c79c44052e1f7bacc140dc5b069a5c407974a5a05050d4a835059
# # FROM  python:3.7-slim-bullseye
# # FROM python:latest
# FROM python:3.11.0a1

# RUN mkdir /tensorflow
# WORKDIR /tensorflow

# COPY extract.py .
# COPY requirements.txt .
# COPY cloud_key.json .docker 
# COPY loop.sh .

# RUN apt-get install bash  -y
# RUN apt-get install curl  -y
# # RUN pip3 install --upgrade setuptools
# # RUN python3 -m pip install --upgrade pip

# RUN pip3 install -r requirements.txt

# CMD ["/tensorflow/loop.sh"]
# # CMD [ "python3", "/tensorflow/extract.py" ]









# FROM nvcr.io/nvidia/l4t-tensorflow:r32.6.1-tf2.5-py3


# RUN mkdir /tensorflow
# WORKDIR /tensorflow

# COPY extract.py .
# COPY requirements.txt .
# COPY cloud_key.json .

# # RUN pip3 install --upgrade setuptools
# # RUN python3 -m pip install --upgrade pip

# RUN pip3 install -r requirements.txt

# # CMD ["/tensorflow/loop.sh"]
# CMD [ "python3", "/tensorflow/extract.py" ]








FROM debian:buster-slim AS build
RUN apt-get update && \
    apt-get install --no-install-suggests --no-install-recommends --yes python3-venv gcc libpython3-dev && \
    python3 -m venv /venv && \
    /venv/bin/pip install --upgrade pip


FROM build AS build-venv
COPY requirements.txt /requirements.txt
RUN /venv/bin/pip install --disable-pip-version-check -r /requirements.txt

#  distroless python image
FROM gcr.io/distroless/python3-debian10
COPY --from=build-venv /venv /venv
COPY . /app
WORKDIR /app
ENTRYPOINT ["/venv/bin/python3", "extract.py"]
