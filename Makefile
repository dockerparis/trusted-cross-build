all: build

build:  qemu-arm-static
	docker build -t trusted-docker-build .

qemu-arm-static:
	cp $(shell which qemu-arm-static) $@

run: build
	docker run -it --rm trusted-docker-build
