REGISTRY=docker.pkg.github.com/elwin
REPOSITORY=heatmap
VERSION=latest

build:
	docker build . -t $(REPOSITORY)

release: build
	docker tag $(REPOSITORY) $(REGISTRY)/$(REPOSITORY)/$(REPOSITORY):$(VERSION)
	docker push $(REGISTRY)/$(REPOSITORY)/$(REPOSITORY):$(VERSION)

test:
	go test ./...