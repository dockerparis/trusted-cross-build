all: build

build:  Dockerfile qemu-arm-static/qemu-arm-static wrapper/wrapper-i386
	docker build -t trusted-docker-build .

wrapper/wrapper-i386:
	make -C wrapper wrapper

qemu-arm-static/qemu-arm-static:
	make -C qemu-arm-static qemu-arm-static

run: build
	docker run -it --rm trusted-docker-build
