FROM ubuntu:latest
RUN apt-get -y update
COPY workshop_demo /opt/cloud/service/
USER root