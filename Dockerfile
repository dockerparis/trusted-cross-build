FROM armbuild/busybox:latest

ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
#RUN /usr/bin/qemu-arm-static /bin/sh -c "echo ':arm:M::\x7fELF\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x28\x00:\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff:/usr/bin/qemu-arm-static:' > /sys/fs/binfmt_misc/register"
RUN /usr/bin/qemu-arm-static /bin/sh -c 'echo Hello World !'
RUN sh -c 'echo Hello World !'
RUN echo Hello World !
