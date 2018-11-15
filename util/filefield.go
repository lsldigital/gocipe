package util

// FileField denotes a field which hold information of files (name) on disk (e.g. uploaded files)
type FileField struct {
	ConfigName     string
	Destination    string
	EntityName     string
	FieldName      string
	SchemaName     string
	ContentBuilder bool
}
