build:
	go build generator.go
	./generator -src=./config/config.go -type=Config -output=config_merge.go 