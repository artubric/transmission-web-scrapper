APP_NAME = tws
VERSION = 0.1.4
COMPILE_NAME = ${APP_NAME}-${VERSION}

BUILD_OS_TARGET = $(shell go env GOOS)
BUILD_ARCH_TARGET = $(shell go env GOARCH)

# build local executable
build-local:
	go build -o bin/${COMPILE_NAME} cmd/transmission-web-scrapper/main.go

# build docker image (for local use)
build-docker-image: set-docker-tag
	docker build -f ./build/Dockerfile . \
		-t ${DOCKER_TAG} \
		--build-arg BUILD_OS_TARGET=${BUILD_OS_TARGET} \
		--build-arg BUILD_ARCH_TARGET=${BUILD_ARCH_TARGET}

# build and save docker image for local use
build-docker-image-local: build-docker-image save-docker-image

# build docker image for Synology DS920+
build-docker-image-synology-920p: set-synology-920p-vars build-docker-image save-docker-image

set-synology-920p-vars:
	$(eval BUILD_OS_TARGET=linux)
	$(eval BUILD_ARCH_TARGET=amd64)

set-docker-tag:
	$(eval DOCKER_TAG=${APP_NAME}:${BUILD_OS_TARGET}-${BUILD_ARCH_TARGET}-${VERSION})

save-docker-image:
	docker save -o ./bin/di-${APP_NAME}-${BUILD_OS_TARGET}-${BUILD_ARCH_TARGET}-${VERSION} ${DOCKER_TAG}
