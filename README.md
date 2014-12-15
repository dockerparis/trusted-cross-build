Trusted cross build
===================

**The author is a jury member, this project is out of the hackathon**

Cross compilation image tools for Docker trusted build (POC)

The goal was to mimic the `binfmt-support` of the Kernel in a Dockerfile

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

But as soon as we have an execve, it breaks ([build](https://registry.hub.docker.com/u/moul/trusted-cross-build/build_id/29996/code/bm8ck9q56mqugabcaxxwcmi/)):

    FROM armbuild/ubuntu:latest
    ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
    ADD wrapper/wrapper-i386 /bin/sh
    RUN date

    Step 0 : FROM armbuild/ubuntu:latest
    Pulling image (latest) from armbuild/ubuntu, endpoint: https://registry-1.docker.io/v1/ 7ae58afd9325
    [...]
    Step 5 : RUN date
    ---> Running in 1abb5e69e50b
    x86_64
    [/bin/sh -c date]
    [/usr/bin/qemu-arm-static /bin/bash -c date]
    [91m/bin/bash: /bin/date: cannot execute binary file: Exec format error
    [0m
    The command [/bin/sh -c date] returned a non-zero code: 126

Second Proof-of-concept ([build](https://registry.hub.docker.com/u/moul/trusted-cross-build/build_id/29996/code/bj3ref9df6a6dmzdeznnkrf/)):

The wrapper looks like `busybox`, we have a symlink pointing to the wrapper for each binaries we will want to run.

The wrapper will look for `argv[0]` and call the original binary prefixed with `qemu-arm-static`.

It works for more cases, but we need to prefix each binaries, this solution can be improved with another wrapper that will list the wanted binaries, then generate all the needed symlinks.

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
    RUN date
    RUN /bin/date
    RUN /bin/bash -c /bin/date
    RUN bash -c /bin/date
    RUN /bin/sh -c /bin/date
    RUN sh -c /bin/date
    # RUN apt-get install -y cowsay # -> failing with: Failed to exec method /usr/lib/apt/methods/http
    CMD ["bash"]

    Step 0 : FROM armbuild/ubuntu:latest
    Pulling image (latest) from armbuild/ubuntu, endpoint: https://registry-1.docker.io/v1/ 7ae58afd9325
    Download complete 7ae58afd9325
    Download complete 7ae58afd9325
    Status: Downloaded newer image for armbuild/ubuntu:latest
    ---> 7ae58afd9325
    Step 1 : ADD qemu-arm-static/qemu-arm-static /usr/bin/qemu-arm-static
    ---> 6ac9a3150cdb
    Removing intermediate container c20f488e0e5f
    Step 2 : ADD wrapper/wrapper-i386 /bin/sh
    ---> 3c51e8c14252
    Removing intermediate container 9246489850a9
    Step 3 : ADD ld_wrapper/ld_wrapper.so /bin/ld_wrapper.so
    ---> 8018ef2e3f38
    Removing intermediate container 198ddb4d1123
    Step 4 : ADD binproxy /binproxy
    ---> 12534e9daba7
    Removing intermediate container 492797f2e2de
    Step 5 : ENV PATH /binproxy
    ---> Running in cc1a8cdea8d3
    ---> 9da10af1a80e
    Removing intermediate container cc1a8cdea8d3
    Step 6 : RUN sh -c 'echo Hello World !'
    ---> Running in bbeb0f5d3baa
    no such file or directory: /bin/wrapper
    ---> 29ad7ad4c922
    Removing intermediate container bbeb0f5d3baa
    Step 7 : RUN echo Hello World !
    ---> Running in 5859091b1fd3
    Input args:  [/bin/sh -c echo Hello World !]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c echo Hello World !]
    Hello World !
    ---> 596b77bc9865
    Removing intermediate container 5859091b1fd3
    Step 8 : RUN echo Hello World !
    ---> Running in 9df4301aeb4e
    Input args:  [/bin/sh -c echo Hello World !]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c echo Hello World !]
    Hello World !
    ---> 5aace60cbe09
    Removing intermediate container 9df4301aeb4e
    Step 9 : RUN echo Hello World !
    ---> Running in c8ad5d2829f7
    Input args:  [/bin/sh -c echo Hello World !]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c echo Hello World !]
    Hello World !
    ---> 92462ed30bdf
    Removing intermediate container c8ad5d2829f7
    Step 10 : RUN date
    ---> Running in 5d698660eb0b
    Input args:  [/bin/sh -c date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c date]
    Input args:  [date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:36:55 UTC 2014
    ---> 41744084d1fd
    Removing intermediate container 5d698660eb0b
    Step 11 : RUN /bin/date
    ---> Running in 081b67fe5025
    Input args:  [/bin/sh -c /bin/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/date]
    Input args:  [/binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:36:57 UTC 2014
    ---> 599c2f87ad48
    Removing intermediate container 081b67fe5025
    Step 12 : RUN /bin/bash -c /bin/date
    ---> Running in 5d99ee0c13df
    Input args:  [/bin/sh -c /bin/bash -c /bin/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/bash -c /binproxy/date]
    Input args:  [/binproxy/bash -c /binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/date]
    Input args:  [/binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:36:59 UTC 2014
    ---> 1e9eaf4d5fde
    Removing intermediate container 5d99ee0c13df
    Step 13 : RUN bash -c /bin/date
    ---> Running in 17b7dd5043c8
    Input args:  [/bin/sh -c bash -c /bin/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c bash -c /binproxy/date]
    Input args:  [bash -c /binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/date]
    Input args:  [/binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:37:01 UTC 2014
    ---> 1a71f497a34f
    Removing intermediate container 17b7dd5043c8
    Step 14 : RUN /bin/sh -c /bin/date
    ---> Running in ba2bce255960
    Input args:  [/bin/sh -c /bin/sh -c /bin/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/sh -c /binproxy/date]
    Input args:  [/binproxy/sh -c /binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/date]
    Input args:  [/binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:37:03 UTC 2014
    ---> fc0b7dc8a8c2
    Removing intermediate container ba2bce255960
    Step 15 : RUN sh -c /bin/date
    ---> Running in 579c2effea5c
    Input args:  [/bin/sh -c sh -c /bin/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c sh -c /binproxy/date]
    Input args:  [sh -c /binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/bash -c /binproxy/date]
    Input args:  [/binproxy/date]
    Output args: [/usr/bin/qemu-arm-static /bin/date]
    Mon Dec 15 07:37:05 UTC 2014
    ---> 36ab144cbf99
    Removing intermediate container 579c2effea5c
    Step 16 : CMD bash
    ---> Running in 8e23dde8cfae
    ---> 0afaa5c20659
    Removing intermediate container 8e23dde8cfae
    Successfully built 0afaa5c20659

And the essential, it runs on an armhf machine

    root@devbox-image-tools-docker-builder:~# uname -a
    Linux devbox-image-tools-docker-builder 3.17.0-119 #1 SMP Thu Nov 20 14:15:44 CET 2014 armv7l armv7l armv7l GNU/Linux
    root@devbox-image-tools-docker-builder:~# docker run -it --rm moul/trusted-cross-build /bin/bash
    root@0a2a778ad62c:/# uname -a
    Linux 0a2a778ad62c 3.17.0-119 #1 SMP Thu Nov 20 14:15:44 CET 2014 armv7l armv7l armv7l GNU/Linux
    root@0a2a778ad62c:/#

Evolution
---------

A third proof-of-concept that uses a dynamic library that will be `LD_PRELOAD`ed and hot prefix the `execve` calls.

So we should have :

- `qemu-arm-static` official binary
- a `wrapper.so` armhf dynamic library
- /bin/sh wrapper: i386 binary that will call the targeted binary with `qemu-arm-static` **and** set the `$LD_PRELOAD` environment variable for all the children processes

I think this solution should be better than the last one; Maybe for the next hackaton :)

---

Better clean and packaging of this hack, something like :

- add a "ADD http://j.mp/install-qemu-arm-hack.tar /hack" just after the `FROM ...`
- add a "RUN /hack/clean" at the end of the Dockerfile

Author
------

- [Manfred Touron](https://github.com/moul)
