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
ENTRYPOINT ["/venv/bin/python3", "validate.py"]




# FROM python:3.7-slim-bullseye


# # RUN apt-get install bash -y

# RUN mkdir /tensorflow
# WORKDIR /tensorflow

# COPY validate.py .
# COPY loop.sh .
# COPY training_dataset_1 .
# COPY requirements.txt .


# RUN pip3 install -r requirements.txt

# # CMD ["/tensorflow/loop.sh"]
# CMD [ "python3", "/tensorflow/validate.py" ]
