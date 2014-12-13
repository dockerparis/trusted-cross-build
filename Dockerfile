FROM armbuild/ubuntu:latest

ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
ADD wrapper/wrapper-i386 /bin/sh

RUN sh -c 'echo Hello World !'
RUN echo Hello World !
RUN date
RUN echo 42

CMD bash