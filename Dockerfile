FROM alpine:3.18

RUN apk add --no-cache git bash openssh

COPY ./workshop_demo /usr/local/bin

ENTRYPOINT ["/usr/local/bin/workshop_demo"]