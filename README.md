# gocipe
Recipes for go


## Usage

```go
type Person struct {
    ID       *int64  `json:"id"       field.name:"id"        field.type:"serial"`
    Password *string `json:"-"        field.name:"password" field.type:"varchar(128)"`
    Alias    *string `json:"alias"    field.name:"alias"     field.type:"varchar(32)" field.nullable:"true"`
    Name     *string `json:"name"     field.name:"name"      field.type:"varchar(255)" field.default:"John Doe"`
    Status   *string `json:"status"   field.name:"status"    field.type:"char(1)"`
}
```

## Tags

Tag name          | Description
------------------|-----------------------------------------------------
json              | Defines JSON key to be encoded
field.name        | Defines DB table column name
field.type        | Defines DB table column type
field.nullable    | "true" (nullable) or "false" (not null). Default: "false"
field.default     | Defines default value of field for database schema
field.filterable  | If field can be used as filter for REST endpoints. "true" / "false". Default: "true"
widget            | Widget to use (vuetify). Format is `Label#Type` or `Label#Type#Options`

## Widget

Type         | Description                                                                | Options
-------------|----------------------------------------------------------------------------|--------
`textfield`  | Text box                                                                   |   
`textarea`   | Text area                                                                  |    
`number`     | Number                                                                     | 
`password`   | Password                                                                   |   
`checkbox`   | Checkbox                                                                   |   
`date`       | Date picker                                                                |      
`select`     | Select with predefined options                                             | `key1|label1;key2|label2;keyn|labeln`             
`select-rel` | Select with options fetched asynchronously from a related entity endpoint  | `endpoint;filtername` where filtername is the field used for label                                                        