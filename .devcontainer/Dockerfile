FROM golang:1.17

ENV DEBIAN_FRONTEND=noninteractive

ENV PATH=$PATH:~/bin

COPY . /dockerbuild/app
RUN /dockerbuild/app/build-dev.sh
