# syntax=docker/dockerfile:1

FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add bash
RUN apk add --no-cache libc6-compat 

COPY KNXDataExposer /
WORKDIR /app
COPY *.yaml .
EXPOSE 12345
ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
CMD [ "/KNXDataExposer" ]