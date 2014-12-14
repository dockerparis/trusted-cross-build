build: ld_wrapper.so Dockerfile

ld_wrapper.so: ld_wrapper.o
	docker build -t gcc-arm .
	docker run -it --rm gcc-arm cat ld_wrapper.so > ld_wrapper.so
