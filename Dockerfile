FROM armbuild/busybox:latest

ADD qemu-arm-static /usr/bin/qemu-arm-static

RUN echo Hello World !


