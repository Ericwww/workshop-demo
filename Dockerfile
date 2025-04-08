FROM --platform=linux/amd64 ubuntu:latest

RUN apt-get install git bash openssh

COPY ./workshop_demo /usr/local/bin/workshop_demo

ENTRYPOINT ["/usr/local/bin/workshop_demo"]

