FROM ubuntu:latest

RUN apt-get -y install gcc-4.8-arm-linux-gnueabihf
ADD ld_wrapper.c /build/
WORKDIR /build/
ENV gcc /usr/bin/arm-linux-gnueabihf-gcc-4.8

RUN $gcc -c -Wall -fpic ld_wrapper.c
RUN $gcc -o ld_wrapper.so -shared ld_wrapper.o
