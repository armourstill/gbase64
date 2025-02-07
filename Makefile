bin:
	go build -o gbase64 -mod vendor github.com/armourstill/gbase64/cmd

vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;



