all:	build

build:  Dockerfile qemu-arm-static/qemu-arm-static ld_wrapper/ld_wrapper.so wrapper/wrapper-i386
	docker build -t trusted-docker-build .

wrapper/wrapper-i386: wrapper/wrapper.go
	make -C wrapper build

ld_wrapper/ld_wrapper.so: ld_wrapper/ld_wrapper.c
	make -C ld_wrapper build

binproxy/sh: binproxy/Makefile
	make -C binproxy build

qemu-arm-static/qemu-arm-static: qemu-arm-static/Dockerfile
	make -C qemu-arm-static build

run:	build
	docker run -it --rm trusted-docker-build /usr/bin/qemu-arm-static /bin/bash

clean:
	make -C wrapper clean
	make -C qemu-arm-static clean
	make -C ld_wrapper clean
	make -C binproxy clean
