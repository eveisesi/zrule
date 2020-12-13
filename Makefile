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

prod:
	git checkout --orphan dist
	git --work-tree .dist add --all
	git --work-tree .dist commit -m "dist"
	git push origin HEAD:dist --force
	git checkout -f main
	git branch -D dist