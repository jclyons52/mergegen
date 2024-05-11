package config

import (
	"time"
)

// generated code, do not modify

type Merger struct {
	MergeCreatedAt func(dst, src *time.Time) error
}

// NewMerger creates a new Merger with optional custom merge functions for external types.
func NewMerger(
	mergeCreatedAt func(dst, src *time.Time) error,
) *Merger {
	return &Merger{
		MergeCreatedAt: mergeCreatedAt,
	}
}

// MergeFeatures merges two Features structs.
func (m *Merger) MergeFeatures(dst, src *Features) error {
	if src.EnableLogging != false {
		dst.EnableLogging = src.EnableLogging
	}
	if src.MaxRetries != 0 {
		dst.MaxRetries = src.MaxRetries
	}
	return nil
}

// MergeClient merges two Client structs.
func (m *Merger) MergeClient(dst, src *Client) error {
	if src.Host != "" {
		dst.Host = src.Host
	}
	if src.Port != 0 {
		dst.Port = src.Port
	}
	return nil
}

// MergeConfig merges two Config structs.
func (m *Merger) MergeConfig(dst, src *Config) error {
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
		if err := m.MergeFeatures(dst.Features, src.Features); err != nil {
			return err
		}
	}
	if err := m.MergeClient(&dst.Client, &src.Client); err != nil {
		return err
	}
	if src.Bar != "" {
		dst.Bar = src.Bar
	}
	if err := m.MergeCreatedAt(&dst.CreatedAt, &src.CreatedAt); err != nil {
		return err
	}
	return nil
}
