all: build

build:  Dockerfile qemu-arm-static/qemu-arm-static wrapper/wrapper
	docker build -t trusted-docker-build .

wrapper/wrapper:
	make -C wrapper wrapper

qemu-arm-static/qemu-arm-static:
	make -C qemu-arm-static qemu-arm-static

run: build
	docker run -it --rm trusted-docker-build
