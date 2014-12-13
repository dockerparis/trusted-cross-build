wrapper-i386: Dockerfile wrapper.go
	docker build -t wrapper .
	docker run wrapper true
	$(eval CID := $(shell docker ps -lq))
	docker cp $(CID):/go/bin tmp
	mv tmp/bin/* .
	rm -rf tmp

