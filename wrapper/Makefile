wrapper-i386: Dockerfile wrapper.go
	docker build -t wrapper .
	docker run wrapper true
	docker cp $(shell docker ps -lq):/go/bin tmp
	mv tmp/bin/* .
	rm -rf tmp

