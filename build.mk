docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u $(DOCKER_USERNAME) --password-stdin

docker-build: docker-login
	docker build --platform linux/amd64 -t muhammed2534/foodtinder:latest -f .docker/Dockerfile .
	docker push muhammed2534/foodtinder:latest