vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;

.PHONY: gobase64
gobase64:
	go build -mod vendor ${FLAGS} -o bin/gobase64 main.go
