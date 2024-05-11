package config

import (
	"github.com/imdario/mergo"
)

// generated code, do not modify

// MergeFeatures merges two Features structs.
func MergeFeatures(dst, src *Features) error {
	if src.EnableLogging != false {
		dst.EnableLogging = src.EnableLogging
	}
	if src.MaxRetries != 0 {
		dst.MaxRetries = src.MaxRetries
	}
	return nil
}

// MergeClient merges two Client structs.
func MergeClient(dst, src *Client) error {
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Port != 0 {
		dst.Port = src.Port
	}
	return nil
}

// MergeConfig merges two Config structs.
func MergeConfig(dst, src *Config) error {
	if src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if src.Timeout != 0 {
		dst.Timeout = src.Timeout
	}
	if src.Features != nil {
		if dst.Features == nil {
			dst.Features = new(Features)
		}
		if err := MergeFeatures(dst.Features, src.Features); err != nil {
			return err
		}
	}
	if err := MergeClient(&dst.Client, &src.Client); err != nil {
		return err
	}
	if src.Bar != "" {
		dst.Bar = src.Bar
	}
	if err := mergo.Merge(&dst.CreatedAt, src.CreatedAt, mergo.WithOverride); err != nil {
		return err
	}
	return nil
}
