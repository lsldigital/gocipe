# gocipe
Recipes for go

**Gocipe** generates, for an entity, the following

- CRUD functions and methods

- database schema

- http server

- REST endpoint functions and methods

## Usage

Group together data in a `struct` to form records (entity). You can
define a struct, for example, in models/person/person.go

```go
type Person struct {
    ID       *int64  `json:"id"     field.name:"id"       field.type:"serial"`
    Password *string `json:"-"      field.name:"password" field.type:"varchar(128)"`
    Alias    *string `json:"alias"  field.name:"alias"    field.type:"varchar(32)"   field.nullable:"true"`
    Name     *string `json:"name"   field.name:"name"     field.type:"varchar(255)"  field.default:"John Doe"`
    Status   *string `json:"status" field.name:"status"   field.type:"char(1)"`
}

//go:generate gocipe crud -file $GOFILE -struct Person
//go:generate gocipe rest -file $GOFILE -struct Person
//go:generate gocipe db   -file $GOFILE -struct Person -output ../../db/persons.sql
```

On the command line:

```bash
$ go generate ./models/person
```

Some packages may have not been imported. Run this:

```bash
$ goimports -w models && gofmt -w models
```

## Tags

Tag name          | Description
------------------|-----------------------------------------------------
json              | Defines JSON key to be encoded
field.name        | Defines DB table column name
field.type        | Defines DB table column type
field.nullable    | "true" (nullable) or "false" (not null). Default: "false"
field.default     | Defines default value of field for database schema
field.filterable  | Whether field can be used as filter for REST endpoints. "true" / "false". Default: "true"

