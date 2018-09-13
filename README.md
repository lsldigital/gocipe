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
field.filterable  | If field can be used as filter for REST endpoints. "true" / "false". Default: "true"
field.mmany       | Format is `thisid#pivotTable#thatid`
widget            | Widget to use (vuetify). Format is `Label#Type` or `Label#Type#Options`

## Command, Arguments & Flags

Go:generate accepts a command, one or more argument and one or more flags. For example, in the following:

```go
//go:generate gocipe crud -file $GOFILE -struct Person
```

- `gocipe` is the command
- `crud` is the argument
- `file` and `struct` are flags

**Notes**
There is no guarantee about imports and formatting. Some imports may be missing and some files may have loose formatting.
Use goimports and go fmt to ensure imports and formatting are correct, example: `goimports -w models/...` and `go fmt -w models/...`

### Arguments

`Gocipe` accepts one mandatory argument.

Argument   | Description
-----------|-------------------------------------
http       | Generates http server boiler plate. Use in `main.go`
crud       | Generates CRUD functions and methods
db         | Generates database schema
rest       | Generates REST endpoint functions and methods
vuetify    | Generates vuetify Edit and List components for entities

The arguments have some flags associated to them.

### common

Flag     | Required | Default | Description
---------|----------|---------|------------------------------------------
file     | Yes      |         | Filename where struct is located
struct   | Yes      |         | Name of the structure to use
v        | No       | false   | Verbose mode. False by default

### crud
Flag     | Required | Default | Description
---------|----------|---------|------------------------------------------
`DELETE` |          |         |
d        | No       | true    | Generate Delete
hd       | No       | false   | Generate Delete pre-execution hook
dh       | No       | false   | Generate Delete post-execution hook
`GET`    |          |         |
g        | No       | true    | Generate Get
hg       | No       | false   | Generate Get pre-execution hook
gh       | No       | false   | Generate Get post-execution hook
`SAVE`   |          |         |
s        | No       | true    | Generate Save
hs       | No       | false   | Generate Save pre-execution hook
sh       | No       | false   | Generate Save post-execution hook
`LIST`   |          |         |
l        | No       | true    | Generate List
hl       | No       | false   | Generate List pre-execution hook
lh       | No       | false   | Generate List post-execution hook

### db
Flag     | Required | Default | Description
---------|----------|---------|------------------------------------------
output   | Yes      |         | SQL filename to write to

```
//go:generate gocipe db -file $GOFILE -struct Person -output $GOPATH/src/github.com/namespace/project/db/person.sql
```

### http

```
//go:generate gocipe http -file $GOFILE
```

### rest

Flag     | Required | Default | Description
---------|----------|---------|------------------------------------------
file     | Yes      |         | Filename where struct is located
struct   | Yes      |         | Name of the structure to use
v        | No       | false   | Verbose mode. False by default
`DELETE` |          |         |
d        | No       | true    | Generate Delete
hd       | No       | false   | Generate Delete pre-execution hook
dh       | No       | false   | Generate Delete post-execution hook
`GET`    |          |         |
g        | No       | true    | Generate Get
hg       | No       | false   | Generate Get pre-execution hook
gh       | No       | false   | Generate Get post-execution hook
`CREATE` |          |         |
c        | No       | true    | Generate Create
hc       | No       | false   | Generate Create pre-execution hook
ch       | No       | false   | Generate Create post-execution hook
`UPDATE` |          |         |
u        | No       | true    | Generate Update
hu       | No       | false   | Generate Update pre-execution hook
uh       | No       | false   | Generate Update post-execution hook
`LIST`   |          |         |
l        | No       | true    | Generate List
hl       | No       | false   | Generate List pre-execution hook
lh       | No       | false   | Generate List post-execution hook

## Widget

All widget tags follow the format `Label#Type(Data)#Options`

Type         | Description                                                               | Data                                          | Options
-------------|---------------------------------------------------------------------------|-----------------------------------------------|--------
`textfield`  | Text box                                                                  |                                               |   
`textarea`   | Text area                                                                 |                                               |    
`number`     | Number                                                                    |                                               | 
`range `     | Range                                                                     |                                               |  
`password`   | Password                                                                  |                                               |   
`checkbox`   | Checkbox                                                                  |                                               |   
`date`       | Date picker                                                               |                                               |      
`select`     | Select with predefined options                                            | `(key1:value1, key2:value2, ... keyn:value1)` |              
`select-rel` | Select with options fetched asynchronously from a related entity endpoint | `(endpoint,filtername)` e.g: `(persons;name)` |                                                        
> If no data, then it must be ommitted, including the `()`


## Update for related entities

Current case: Save triggered on Entity A
| Relationship Types | Entity A  | Entity B   | Notes                                            |
|--------------------|-----------|------------|--------------------------------------------------|
| One-One            | Saved     | Not Saved* | * Save foreign key value of entity A  (i.e UUID) |
| One-Many           | Saved     | Not Saved  |                                                  |
| Many-One           | Saved*    | Not Saved  | * Save foreign key value of entity B (i.e UUID)  |
| Many-Many          | Saved     | Saved      |                                                  |
| Many-Many-Owner    | Saved     | Saved*     | * Only in pivot table                            |
| Many-Many-Inverse  | Saved     | Not Saved  |                                                  |

### Notes
> * The owning side of a `one-to-one` association is the entity with the table containing the foreign key.
> * `One-to-many` is always the inverse side of a bidirectional association.
> * `Many-to-one` is always the owning side of a bidirectional association.
> * You can pick the owning side of a `many-to-many` association yourself.