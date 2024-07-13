.PHONY: docker-build-server
docker-build-server:
	@docker build -f _docker/server.Dockerfile -t kube-registry:5000/boardgames:latest --cache-from kube-registry:5000/boardgames:latest .

.PHONY: docker-push-server
docker-push-server:
	@docker push kube-registry:5000/boardgames:latest

.PHONY: k8s-apply
k8s-apply:
	@kubectl apply -f _infrastructure/kube/manifests/server.yml

.PHONY: k8s-deploy
k8s-deploy: docker-build-server docker-push-server k8s-apply
