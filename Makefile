.PHONY: docker-push
docker-push:
	templ generate
	docker build -t m4tthewde-paste .
	docker tag m4tthewde-paste:latest 211018008663.dkr.ecr.us-east-1.amazonaws.com/m4tthewde-paste:latest
	docker push 211018008663.dkr.ecr.us-east-1.amazonaws.com/m4tthewde-paste:latest
