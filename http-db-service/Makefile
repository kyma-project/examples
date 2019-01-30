APP_NAME = http-db-service
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve: 
	dep ensure --vendor-only -v

validate:
	curl https://raw.githubusercontent.com/alecthomas/gometalinter/master/scripts/install.sh | sh -s v2.0.8
	${GOPATH}/bin/gometalinter --vendor --deadline=2m --disable-all --enable=vet ./...
	
build:
	go generate ./...
	CGO_ENABLED=0 go build -o ./bin/app $(buildpath)

test-report:
	2>&1 go test -v ./... | go2xunit -fail -output unit-tests.xml

.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest .

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)

.PHONY: ci-pr
ci-pr: resolve validate build-image push-image

.PHONY: ci-master
ci-master: resolve validate build-image push-image