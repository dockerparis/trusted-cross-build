qemu-arm-static: Dockerfile
	docker build -t qemu-arm-static .
	docker run qemu-arm-static true
	docker cp $(shell docker ps -lq):/usr/bin/qemu-arm-static tmp
	touch tmp/qemu-arm-static
	mv tmp/qemu-arm-static qemu-arm-static
	rmdir tmp

