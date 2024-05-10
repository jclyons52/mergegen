build:
	go build .
run:
	go run . -src=./config/config.go -type=Config -output=config_merge.go
test:
	go test -v -cover ./...