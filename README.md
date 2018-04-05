# gocipe
Recipes for go


### StructureInfo

StructureInfo represents a target structure in a file to be used for generation

```go
type StructureInfo struct {
	Package   string
	Name      string
	TableName string
	Fields    []FieldInfo
}
```

### FieldInfo

FieldInfo represents information about a field

```go
type FieldInfo struct {
	Name     string				// field.name
	Property string				// GO struct fields (ID, Authcode, ...)
	Type     string				// GO basic value types (int64, string, ...) or custom types
	DBType   string				// field.type
	Tags     reflect.StructTag  // GO struct field tags (between ``)
}
```

Example of fields (in a struct)

```go
type Struct struct {
    ID       *int64  `json:"id"       field.name:"id"        field.type:"serial"`
    Authcode *string `json:"-"        field.name:"auth_code" field.type:"varchar(128)"`
    Alias    *string `json:"alias"    field.name:"alias"     field.type:"varchar(32)"`
    Name     *string `json:"name"     field.name:"name"      field.type:"varchar(255)"`
    Callback *string `json:"callback" field.name:"callback"  field.type:"varchar(255)"`
    Status   *string `json:"status"   field.name:"status"    field.type:"char(1)"`
}
```

**Tags**

Tag name          | Description
------------------|-----------------------------------------------------
json              | Defines JSON key to be encoded
field.name        | Defines DB table column name
field.type        | Defines DB table column type
field.nullable    | "true" (nullable) or "false" (not null). Default: "false"
field.default     | Defines default value of field
