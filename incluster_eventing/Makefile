.DEFAULT_GOAL := render
render:
	(cd ./src ; sh render-manifests.sh)
deploy:
	kubectl apply -f ./k8s-resources