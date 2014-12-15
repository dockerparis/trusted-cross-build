FROM armbuild/ubuntu:latest

# PREPARE IMAGE
ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
ADD wrapper/wrapper-i386 /bin/sh
ADD ld_wrapper/ld_wrapper.so /bin/ld_wrapper.so
ADD binproxy /binproxy
ENV PATH /binproxy
#ENV LD_PRELOAD /bin/ld_wrapper.so

# STANDARD COMMANDS
RUN sh -c 'echo Hello World !'
RUN echo Hello World !
RUN echo Hello World !
RUN echo Hello World !
#RUN /bin/wrapper /bin/date
RUN date
#RUN /bin/date
#RUN /bin/wrapper /bin/date
#RUN date
#RUN /bin/sh -c /bin/date
#RUN apt-get install -y cowsay
#RUN date
#CMD ["/bin/sh", "-c", "ls -la"]
CMD ["/usr/bin/qemu-arm-static", "/bin/bash"]

# CLEAN IMAGE
