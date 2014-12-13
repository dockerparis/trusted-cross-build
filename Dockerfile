FROM armbuild/busybox:latest
#FROM ubuntu:latest
#RUN apt-get update && apt-get -y install qemu-user-static

ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
ADD wrapper /bin/sh

#ADD qemu-arm-static/qemu-arm-static /bin/sh

#RUN /usr/bin/qemu-arm-static /bin/sh -c "echo ':arm:M::\x7fELF\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x28\x00:\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff:/usr/bin/qemu-arm-static:' > /sys/fs/binfmt_misc/register"
#RUN /usr/bin/qemu-arm-static /bin/sh -c 'echo Hello World !'
#RUN sh -c 'echo Hello World !'
#RUN echo Hello World !

RUN date
RUN echo 42

#RUN /usr/bin/qemu-arm-static /bin/sh -c 'uname'
#RUN /usr/bin/qemu-arm-static /bin/sh -c 'which uname'
#RUN /usr/bin/qemu-arm-static /bin/uname
