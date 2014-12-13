all: build

build:  qemu-arm-static wrapper
	docker build -t trusted-docker-build .

wrapper: wrapper.go
	go build $<

qemu-arm-static:
	cp $(shell which qemu-arm-static) $@

run: build
	docker run -it --rm trusted-docker-build
