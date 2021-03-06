build:
	go build -o zrule-api cmd/zrule/*.go

serve: build
	./zrule-api serve

processor: build
	./zrule-api processor

listener: build
	./zrule-api listener

dispatcher: build
	./zrule-api dispatcher

docker: dockerbuild
dockerbuild:
	docker build . -t zrule:latest


dockercomp: dockercompup dockercomplogs
dockercompup:
	docker-compose up -d --remove-orphans

dockercomplogs:
	docker-compose logs -f serve

dockercompdown:
	docker-compose down
