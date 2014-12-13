wrapper-i386: Dockerfile wrapper.go
	docker build -t wrapper .
	docker run --name=wrapper-tmp wrapper true
	docker cp wrapper-tmp:/go/bin tmp
	docker rm wrapper-tmp
	mv tmp/bin/* .
	rm -rf tmp

