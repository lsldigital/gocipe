package util

// FileField denotes a field which hold information of files (name) on disk (e.g. uploaded files)
type FileField struct {
	ConfigName     string
	Destination    string
	EntityName     string
	FieldName      string
	ContentBuilder bool
}

// GetSchemaName returns field name in snake case format
func (f *FileField) GetSchemaName() string {
	return ToSnakeCase(f.FieldName)
}
