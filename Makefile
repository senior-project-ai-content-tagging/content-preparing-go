image-name = content-preparing
tag = master

run:
	go run ./cmd/main.go

docker-build:
	docker build -t ${image-name}:${tag} .

docker-tag:
	docker tag ${image-name}:${tag} asia.gcr.io/senior-project-364818/${image-name}:${tag}

publish: docker-build docker-tag
	docker push asia.gcr.io/senior-project-364818/${image-name}:${tag}