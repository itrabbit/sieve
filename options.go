package sieve

import "reflect"

type Options interface {
	HasScopes() bool
	HasExportKeys() bool
	HasExclusions() bool

	Scopes() []string
	ExportKeys() []string
	CheckByExclusions(v ...reflect.Value) bool
}

type OptionsBuilder interface {
	Options
	SetScopes(scopes ...string) OptionsBuilder
	SetExportKeys(keys ...string) OptionsBuilder
	PutExclusionStrategy(s ExclusionStrategy) OptionsBuilder
	SetIsAnyExclusion(v bool) OptionsBuilder
}

type options struct {
	scopes         []string            // s
	exportKeys     []string            // k
	exclusions     []ExclusionStrategy // e+ (ef, ev)
	isAnyExclusion bool                // e.any
}

func (o options) HasScopes() bool {
	return o.scopes != nil && len(o.scopes) > 0
}

func (o options) HasExportKeys() bool {
	return o.exportKeys != nil && len(o.exportKeys) > 0
}
func (o options) HasExclusions() bool {
	return o.exclusions != nil && len(o.exclusions) > 0
}

func (o options) Scopes() []string {
	return o.scopes
}

func (o options) ExportKeys() []string {
	return o.exportKeys
}

func (o options) CheckByExclusions(v ...reflect.Value) bool {
	if o.exclusions != nil {
		for _, exclusion := range o.exclusions {
			if exclusion.Check(v...) {
				if o.isAnyExclusion {
					return true
				}
				continue
			} else if !o.isAnyExclusion {
				return false
			}
		}
		return !o.isAnyExclusion
	}
	return false
}

func (o *options) SetScopes(scopes ...string) OptionsBuilder {
	o.scopes = scopes
	return o
}

func (o *options) SetExportKeys(keys ...string) OptionsBuilder {
	o.exportKeys = keys
	return o
}

func (o *options) PutExclusionStrategy(s ExclusionStrategy) OptionsBuilder {
	if o.exclusions == nil {
		o.exclusions = make([]ExclusionStrategy, 0, 0)
	}
	o.exclusions = append(o.exclusions, s)
	return o
}

func (o *options) SetIsAnyExclusion(v bool) OptionsBuilder {
	o.isAnyExclusion = v
	return o
}

func BuildOptions(scopes []string, exportKeys []string) OptionsBuilder {
	return &options{
		scopes:     scopes,
		exportKeys: exportKeys,
	}
}
