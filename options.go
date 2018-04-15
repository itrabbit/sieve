package sieve

type Options interface {
	HasScopes() bool
	HasExportKeys() bool
	HasExcludeEqualField() bool

	Scopes() []string
	ExportKeys() []string
	ExcludeEqualField() string
}

type OptionsBuilder interface {
	Options
	SetScopes(scopes ...string) OptionsBuilder
	SetExportKeys(keys ...string) OptionsBuilder
	SetExcludeEqualField(fieldName string) OptionsBuilder
}

type options struct {
	scopes            []string // s
	exportKeys        []string // ek
	excludeEqualField string   // eef
}

func (o options) HasScopes() bool {
	return o.scopes != nil && len(o.scopes) > 0
}

func (o options) HasExportKeys() bool {
	return o.exportKeys != nil && len(o.exportKeys) > 0
}

func (o options) HasExcludeEqualField() bool {
	return len(o.excludeEqualField) > 0
}

func (o options) Scopes() []string {
	return o.scopes
}

func (o options) ExportKeys() []string {
	return o.exportKeys
}

func (o options) ExcludeEqualField() string {
	return o.excludeEqualField
}

func (o *options) SetScopes(scopes ...string) OptionsBuilder {
	o.scopes = scopes
	return o
}

func (o *options) SetExportKeys(keys ...string) OptionsBuilder {
	o.exportKeys = keys
	return o
}

func (o *options) SetExcludeEqualField(fieldName string) OptionsBuilder {
	o.excludeEqualField = fieldName
	return o
}

func EmptyOptions() OptionsBuilder {
	return &options{}
}

func BuildOptions(scopes []string, exportKeys []string) OptionsBuilder {
	return &options{
		scopes:     scopes,
		exportKeys: exportKeys,
	}
}
