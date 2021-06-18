ARG ARCH
FROM ${ARCH}/ubuntu:20.04

WORKDIR /opt

COPY requirements.txt /tmp

RUN apt-get update && \
    apt-get install -y --no-install-recommends python3-dev python3-pip && \
    python3 -m pip install -r /tmp/requirements.txt && \
    rm /tmp/requirements.txt 

COPY landing-page /opt/landing-page
COPY html /opt/html
COPY assets /opt/assets

CMD python3 /opt/landing-page
