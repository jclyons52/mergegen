build:
	go build -o /usr/local/bin/mergegen
	chmod +x /usr/local/bin/mergegen
run:
	go run . -src=./config/config.go -type=Config -output=config_merge.go
test:
	go test -v -cover ./...