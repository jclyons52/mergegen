package config

import "time"

//go:generate ./generator -src=./config.go -type=Config -output=config_merge.go

type Features struct {
	EnableLogging bool
	MaxRetries    int
}

type Client struct {
	Host string
	Port int
}

type Config struct {
	APIKey    string
	Timeout   int
	Features  *Features
	Client    Client
	Bar       string
	CreatedAt time.Time
}
