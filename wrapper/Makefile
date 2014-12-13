wrapper: Dockerfile wrapper.go
	docker build -t wrapper .
	docker run wrapper true
	docker cp $(shell docker ps -lq):/go/wrapper tmp
	mv tmp/wrapper wrapper
	rmdir tmp

