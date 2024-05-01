package main

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
    if src.Bar != "" {
        dst.Bar = src.Bar
    }
}

