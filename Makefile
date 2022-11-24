
PROJECT_NAME=hackathon
BUILD_VERSION=1.1.0

DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

lambda:
	GOARCH=amd64 GOOS=linux go build -o ./.build/main ./main.go
	zip -jrm ./.build/main.zip ./.build/main

build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin main.go

compose_dev: docker
	cd deploy && BUILD_VERSION=$(BUILD_VERSION) docker compose up --build --force-recreate -d

docker_prebuild: build
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deploy/$(PROJECT_NAME).bin; \

docker_build:
	cd deploy; \
	docker build --rm -t $(DOCKER_IMAGE) .;

docker_postbuild:
	cd deploy; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\

docker: docker_prebuild docker_build docker_postbuild
