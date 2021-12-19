DOCKER_IMAGE?=go-monolith/go-monolith
DOCKER_TAG?=release
DOCKER_IMAGE_TAG?=test
DESTDIR?=$(shell pwd)

.PHONY: docker-image
docker-image: static
	cp ./main contrib/docker/go-monolith.$$(uname -m)
	cp ./configs/sqlite.yml contrib/docker/go-monolith.yml
	docker build -t ${DOCKER_IMAGE}:${DOCKER_IMAGE_TAG} --build-arg ARCH=$$(uname -m) -f contrib/docker/Dockerfile contrib/docker/