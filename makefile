DOCKER_HUB_USERNAME := andreistefanciprian
IMAGE_NAME := stock-ticker-watcher
DOCKER_IMAGE_NAME := $(DOCKER_HUB_USERNAME)/$(IMAGE_NAME)

build:
	docker build -t $(DOCKER_IMAGE_NAME) . -f infra/Dockerfile
	docker image push $(DOCKER_IMAGE_NAME)

test:
	go test  ./... -v