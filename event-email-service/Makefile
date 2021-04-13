APP_NAME = event-email-service
IMG = $(DOCKER_PUSH_REPOSITORY)$(DOCKER_PUSH_DIRECTORY)/$(APP_NAME)
TAG = $(DOCKER_TAG)

resolve:
	npm install --no-optional

validate:
	eslint -c .eslinterrc ./src

test:
	npm run-script test

scan:
	sed -i -e 's/APIKEY/$(API_KEY)/g; s/USERKEY/$(USER_KEY)/g' /whitesource.config.json && cp /whitesource.config.json ./whitesource.config.json
	whitesource run
	
.PHONY: build-image
build-image:
	docker build -t $(APP_NAME):latest .

.PHONY: push-image
push-image:
	docker tag $(APP_NAME) $(IMG):$(TAG)
	docker push $(IMG):$(TAG)

.PHONY: ci-pr
ci-pr: resolve validate test build-image push-image

.PHONY: ci-main
ci-main: resolve validate test build-image push-image
