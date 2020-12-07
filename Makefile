build:
	go build -o zrule cmd/zrule/*.go

http: build
	./zrule http

processor: build
	./zrule processor

listener: build
	./zrule listener

dispatcher: build
	./zrule dispatcher
