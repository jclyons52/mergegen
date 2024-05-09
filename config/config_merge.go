package config

// generated code, do not modify

// MergeFeatures merges two Features structs.
func MergeFeatures(dst, src *Features) {
	if src.EnableLogging != false {
		dst.EnableLogging = src.EnableLogging
	}
	if src.MaxRetries != 0 {
		dst.MaxRetries = src.MaxRetries
	}
}

// MergeClient merges two Client structs.
func MergeClient(dst, src *Client) {
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Port != 0 {
		dst.Port = src.Port
	}
}

// MergeConfig merges two Config structs.
func MergeConfig(dst, src *Config) {
	if src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if src.Timeout != 0 {
		dst.Timeout = src.Timeout
	}
	if dst.Features == nil {
		dst.Features = new(Features)
	}
	MergeFeatures(dst.Features, src.Features)
	if src.Client != nil {
		dst.Client = src.Client
	}
	if src.Bar != "" {
		dst.Bar = src.Bar
	}
}
