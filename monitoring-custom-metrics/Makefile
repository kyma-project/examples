APP_NAME = monitoring-custom-metrics
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve: 
	dep ensure -vendor-only -v

validate:
	go build -o golint-vendored ./vendor/github.com/golang/lint/golint
	./golint-vendored
	rm golint-vendored

build:
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
ci-pr: resolve validate build build-image push-image

.PHONY: ci-master
ci-master: resolve validate build build-image push-image

clean:
	rm -f ./bin/app