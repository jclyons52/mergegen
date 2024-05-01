package main

//go:generate ./generator -src=./config.go -type=Config -output=config_merge.go

type Features struct {
	EnableLogging bool
	MaxRetries    int
}

type Config struct {
	APIKey   string
	Timeout  int
	Features *Features
	Bar      string
}
