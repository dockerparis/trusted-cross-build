Trusted cross build
===================

Cross compilation image tools for Docker trusted build (POC)

Context
-------

Docker images are built for one architecture

    # Local architecture is x86_64
    root@fwrz:~# uname -m
    x86_64
    # We can run x86_64 docker images
    root@fwrz:~# docker run -it --rm busybox uname -m
    x86_64
    # But we cannot run armhf docker images
    root@fwrz:~# docker run -it --rm armbuild/busybox uname -m
    exec format error2014/12/14 20:57:21 Error response from daemon: Cannot start container 37073f0bd91ff94ce670114e9cb2eeef69ee830452ea9712f3c0e2365ec4c0a7: exec format error

By using qemu, we can run binaries built for other architectures

    # We Mount bind the local qemu-arm-static binary and setup the entrypoint, we can now run the armhf image
    root@fwrz:~# apt-get install -y qemu-user-static
    root@fwrz:~# docker run -it --rm -v $(which qemu-arm-static):/usr/bin/qemu-arm-static --rm armbuild/busybox /usr/bin/qemu-arm-static /bin/uname -m
    armv7l

Problem, volumes and custom entrypoints only works for `docker run`, but not for `docker build`.

    root@fwrz:~# cat <<EOF | docker build -t test -
    > FROM armbuild/busybox
    > RUN uname -m
    > EOF
    Sending build context to Docker daemon 2.048 kB
    Sending build context to Docker daemon
    Step 0 : FROM armbuild/busybox
    ---> d91e5575e4cc
    Step 1 : RUN uname -m
    ---> Running in f621d609c524
    exec format error2014/12/14 21:05:20 exec format error

binfmt-support = <3

    root@fwrz:~# apt-get install binfmt-support
    root@fwrz:~# echo ':arm:M::\x7fELF\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x28\x00:\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff:/usr/bin/qemu-arm-static:' > /proc/sys/fs/binfmt_misc/register
    # We can now run the armhf image, by mounting the qemu-arm-static binary but without force the entrypoint, the kernel binfmt-support will take care of this for us
    root@fwrz:~# docker run -it --rm -v $(which qemu-arm-static):/usr/bin/qemu-arm-static --rm armbuild/busybox /bin/uname -m
    armv7l
    # It also works in the build time
    root@fwrz:~# cat Dockerfile
    FROM armbuild/ubuntu:latest
    ADD qemu-arm-static /usr/bin/
    RUN uname -m
    root@fwrz:~# docker build -t toto .
    Sending build context to Docker daemon 29.96 MB
    Sending build context to Docker daemon
    Step 0 : FROM armbuild/ubuntu:latest
    ---> 7ae58afd9325
    Step 1 : ADD qemu-arm-static /usr/bin/
    ---> 4779d849a4dc
    Removing intermediate container 4198cf9a84ed
    Step 2 : RUN uname -m
    ---> Running in 98a74cf8c666
    armv7l
    ---> f68841337716
    Removing intermediate container 98a74cf8c666
    Successfully built f68841337716

Problem, the trusted-build server do not have binfmt-support

Solution
--------

A Dockerfile contains rules that works directly on the filesystem (`FROM`, `ADD`, `COPY`, ...).
The `RUN` command will try to run a shell script line in a container context, it will do this: `docker run $PARENT_CID /bin/sh -c "the shell script line"`.

First proof-of-concept:
- Add `qemu-arm-static` in the container
- Replace the `/bin/sh` in the container (ADD) with a wrapper that will call "qemu-arm-static bash -c $@"

It works for basic commands, but as soon as the wrapped binary do an `execve` syscall, the new child process won't be wrapped.

It worked ([build](https://registry.hub.docker.com/u/moul/trusted-cross-build/build_id/29996/code/biomzg2rvphqzd6yygsw3th/)), There is some debug in the wrapper to see the command translation :

    FROM armbuild/ubuntu:latest
    ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
    ADD wrapper/wrapper-i386 /bin/sh
    RUN sh -c 'echo Hello World !'
    RUN echo Hello World !

    Step 0 : FROM armbuild/ubuntu:latest
    Pulling image (latest) from armbuild/ubuntu, endpoint: https://registry-1.docker.io/v1/ 7ae58afd9325
    Download complete 7ae58afd9325
    Download complete 7ae58afd9325
    Status: Downloaded newer image for armbuild/ubuntu:latest
    ---> 7ae58afd9325
    Step 1 : ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
    ---> 42473d88d32c
    Removing intermediate container 0f3bceb0d97d
    Step 2 : ADD wrapper/wrapper-i386 /bin/sh
    ---> 282501415bb5
    Removing intermediate container 8b6e49b9dfef
    Step 3 : RUN sh -c 'echo Hello World !'
    ---> Running in 3c9ef504d119
    x86_64
    [/bin/sh -c sh -c 'echo Hello World !']
    [/usr/bin/qemu-arm-static /bin/bash -c sh -c 'echo Hello World !']
    x86_64
    [sh -c echo Hello World !]
    [/usr/bin/qemu-arm-static /bin/bash -c echo Hello World !]
    Hello World !
    ---> c78a80ae1416
    Removing intermediate container 3c9ef504d119
    Step 4 : RUN echo Hello World !
    ---> Running in f1598a5e6d6a
    x86_64
    [/bin/sh -c echo Hello World !]
    [/usr/bin/qemu-arm-static /bin/bash -c echo Hello World !]
    Hello World !
    ---> 1f7b6cf626b5


Author
------

- [Manfred Touron](https://github.com/moul)
