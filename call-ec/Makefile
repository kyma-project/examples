APP_NAME = call-ec-web-ui
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve-hybris:
	curl -u $(USERNAME):$(PASSWORD) https://repository.hybris.com/api/npm/auth >> /.npmrc && cp /.npmrc ./.npmrc 
	npm install --no-optional

resolve:
	npm install --no-optional

validate:
	eslint -c .eslinterrc ./web-ui/src

.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest web-ui

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)

.PHONY: ci-pr
ci-pr: resolve validate build-image push-image

.PHONY: ci-master
ci-master: resolve validate build-image push-image