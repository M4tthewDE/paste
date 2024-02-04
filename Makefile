.PHONY: docker-push
docker-push:
	templ generate
	docker build -t m4tthewde-paste .
	docker tag m4tthewde-paste:latest 211018008663.dkr.ecr.us-east-1.amazonaws.com/m4tthewde-paste:latest
	docker push 211018008663.dkr.ecr.us-east-1.amazonaws.com/m4tthewde-paste:latest

.PHONY: docker-login
docker-login:
	@./with-env.sh aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 211018008663.dkr.ecr.us-east-1.amazonaws.com
