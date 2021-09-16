APP_NAME = monitoring-custom-metrics
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve: 
	dep ensure -vendor-only -v

build:
	go mod vendor
	go generate ./...
	CGO_ENABLED=0 go build -o ./bin/app $(buildpath)

.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest .

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)

.PHONY: ci-pr
ci-pr: resolve build build-image push-image

.PHONY: ci-main
ci-main: resolve build build-image push-image

clean:
	rm -f ./bin/app
