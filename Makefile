all: build

build:  Dockerfile qemu-arm-static/qemu-arm-static ld_wrapper/ld_wrapper.so wrapper/wrapper-i386
	docker build -t trusted-docker-build .

wrapper/wrapper-i386: wrapper/wrapper.go
	make -C wrapper

ld_wrapper/ld_wrapper.so: ld_wrapper/ld_wrapper.c
	make -C ld_wrapper

qemu-arm-static/qemu-arm-static: qemu-arm-static/Dockerfile
	make -C qemu-arm-static

run: build
	docker run -it --rm trusted-docker-build
