REGISTRY=ghcr.io/elwin
REPOSITORY=heatmap
NAME=heatmap-ui
VERSION=latest

build:
	docker build . -t $(NAME)

release: build
	docker tag $(NAME) $(REGISTRY)/$(REPOSITORY)/$(NAME):$(VERSION)
	docker push $(REGISTRY)/$(REPOSITORY)/$(NAME):$(VERSION)
