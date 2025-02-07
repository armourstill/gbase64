bin:
	go build -o gbase64 github.com/armourstill/gbase64

vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;
