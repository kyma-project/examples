APP_NAME = http-db-service-acceptance-tests
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve: 
	dep ensure -vendor-only -v

validate:
	go build -o golint-vendored ./vendor/golang.org/x/lint/golint
	./golint-vendored
	rm golint-vendored

.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest .

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)

.PHONY: ci-pr
ci-pr: resolve validate build-image push-image

.PHONY: ci-main
ci-main: resolve validate build-image push-image
