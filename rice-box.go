package main

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "bootstrap.go.tmpl",
		FileModTime: time.Unix(1526036344, 0),
		Content:     string("package app\n\nimport (\n\t\"database/sql\"\n\t\"log\"\n\t\"os\"\n\n\t\"github.com/joho/godotenv\"\n\t// Load database driver\n\t_ \"github.com/lib/pq\"\n)\n\nconst (\n\t//EnvironmentProd represents production environment\n\tEnvironmentProd = \"PROD\"\n\n\t//EnvironmentDev represents development environment\n\tEnvironmentDev  = \"DEV\"\n)\n\nvar (\n\t// bootstrapped is a flag to prevent multiple bootstrapping\n\tbootstrapped = false\n\n\t// Env indicates in which environment (prod / dev) the application is running\n\tEnv string\n\t{{range .Bootstrap.Settings}}{{if .Public}}\n\t// {{.Name}} {{.Description}}\n\t{{.Name}} string\n\t{{end}}{{end}}\n)\n\n// Config represents application configuration loaded during bootstrap\ntype Config struct {\n\t{{if not .Bootstrap.NoDB}}DB  *sql.DB{{end}}\n\t{{if .HTTP.Generate}}HTTPPort string{{end}}\n\t{{range .Bootstrap.Settings}}{{if not .Public}}\n\t// {{.Name}} {{.Description}}\n\t{{.Name}} string\n\t{{end}}{{end}}\n}\n\n// Bootstrap loads environment variables and initializes the application\nfunc Bootstrap() *Config {\n\tvar config Config\n\n\tif bootstrapped {\n\t\treturn nil\n\t}\n\n\tgodotenv.Load()\n\n\tEnv = os.Getenv(\"ENV\")\n\tif Env == \"\" {\n\t\tEnv = EnvironmentProd\n\t}\n\n\t{{if not .Bootstrap.NoDB}}\n\tdsn := os.Getenv(\"DSN\")\n\tif dsn == \"\" {\n\t\tlog.Fatal(\"Environment variable DSN must be defined. Example: postgres://user:pass@host/db?sslmode=disable\")\n\t}\n\n\tvar err error\n\tconfig.DB, err = sql.Open(\"postgres\", dsn)\n\tif err == nil {\n\t\tlog.Println(\"Connected to database successfully.\")\n\t} else if (Env == EnvironmentDev) {\n\t\tlog.Println(\"Database connection failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database connection failed: \", err)\n\t}\n\n\terr = config.DB.Ping()\n\tif err == nil {\n\t\tlog.Println(\"Pinged database successfully.\")\n\t} else if (Env == EnvironmentDev) {\n\t\tlog.Println(\"Database ping failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database ping failed: \", err)\n\t}\n\t{{end}}\n\n\t{{if .HTTP.Generate}}\n\tconfig.HTTPPort = os.Getenv(\"HTTP_PORT\")\n\tif config.HTTPPort == \"\" {\n\t\tconfig.HTTPPort = \"{{.HTTP.Port}}\"\n\t}\n\t{{end}}\n\n\t{{range .Bootstrap.Settings}}{{if not .Public}}\n\tconfig.{{.Name}} = os.Getenv(\"{{upper (snake .Name)}}\")\n\tif config.{{.Name}} == \"\" {\n\t\tlog.Fatal(\"Environment variable {{upper (snake .Name)}} ({{.Description}}) must be defined.\")\n\t}\n\t{{end}}{{end}}\n\n\t{{range .Bootstrap.Settings}}{{if .Public}}\n\t{{.Name}} = os.Getenv(\"{{upper (snake .Name)}}\")\n\tif {{.Name}} == \"\" {\n\t\tlog.Fatal(\"Environment variable {{upper (snake .Name)}} ({{.Description}}) must be defined.\")\n\t}\n\t{{end}}{{end}}\n\n\tos.Clearenv() //prevent non-authorized access\n\n\treturn &config\n}"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "bootstrap_env.tmpl",
		FileModTime: time.Unix(1525884138, 0),
		Content:     string("# The following must be defined as well: ENV{{if not .Bootstrap.NoDB}}, DSN{{end}}{{if .HTTP.Generate}}, HTTP_PORT{{end}}\n{{range .Bootstrap.Settings}}{{upper (snake .Name)}} = \"{{.Description}}\"\n{{end}}"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "crud/crud.go.tmpl",
		FileModTime: time.Unix(1526937997, 0),
		Content:     string("package {{.Package}}\n\nimport (\n\t\"context\"\n\t\"database/sql\"\n\t{{range .Imports}}{{.}}{{end}}\n)\n\nvar db *sql.DB\n\n// Inject allows injection of services into the package\nfunc Inject(database *sql.DB) {\n\tdb = database\n}\n\n{{.Get}}\n{{.List}}\n{{.DeleteSingle}}\n{{.DeleteMany}}\n{{.Save}}\n{{.Insert}}\n{{.Update}}\n{{.Merge}}\n{{range .SaveRelated}}{{.}}{{end}}\n{{range .LoadRelated}}{{.}}{{end}}"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "crud/filters.go.tmpl",
		FileModTime: time.Unix(1526938273, 0),
		Content:     string("package models\n\n//ListFilter represents a filter to apply during listing (crud)\ntype ListFilter struct {\n\tField     string\n\tOperation string\n\tValue     interface{}\n}"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "crud/hooks.go.tmpl",
		FileModTime: time.Unix(1526841119, 0),
		Content:     string("package {{.Package}}\n\nimport (\n\t\"database/sql\"\n)\n\n{{if .Hooks.PreRead}}\nfunc crudPreGet(id {{pkeyPropertyType .Entity.PrimaryKey}}) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PostRead}}\nfunc crudPostGet(entity *{{.Entity.Name}}) error {\n\treturn nil\n}\n{{end}}\n\n{{if .Hooks.PreList}}\nfunc crudPreList(filters []models.ListFilter) ([]models.ListFilter, error) {\n\treturn filters, nil\n}\n{{end}}\n{{if .Hooks.PostList}}\nfunc crudPostList(list []*{{.Entity.Name}}) ([]*{{.Entity.Name}}, error) {\n\treturn list, nil\n}\n{{end}}\n\n{{if .Hooks.PreDelete}}\nfunc crudPreDelete(id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PostDelete}}\nfunc crudPostDelete(id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n\n\n{{if .Hooks.PreSave }}\nfunc crudPreSave(op string, entity *{{.Entity.Name}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PreSave }}\nfunc crudPostSave(op string, entity *{{.Entity.Name}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n\n"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "crud/partials/delete_many.go.tmpl",
		FileModTime: time.Unix(1526939097, 0),
		Content:     string("// Delete deletes many {{.EntityName}} records from database using filter\nfunc Delete(ctx context.Context, filters []models.ListFilter, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr      error\n\t\tstmt     *sql.Stmt\n\t\tsegments []string\n\t\tvalues   []interface{}\n\t)\n\n\tif tx == nil {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\ttx, err = db.Begin()\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t}\n\t}\n\n\tquery := \"DELETE FROM {{.Table}}\"\n\t{{if .HasPreHook}}\n    if filters, err = crudPreDeleteMany(ctx, filters); err != nil {\n\t\treturn err\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tfor i, filter := range filters {\n\t\tsegments = append(segments, filter.Field+\" \"+filter.Operation+\" $\"+strconv.Itoa(i+1))\n\t\tvalues = append(values, filter.Value)\n\t}\n\n\tif len(segments) != 0 {\n\t\tquery += \" WHERE \" + strings.Join(segments, \" AND \")\n\t}\n\n\tstmt, err = db.Prepare(query)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\t_, err = stmt.Exec(values...)\n\t\tif err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif err = crudPostDeleteMany(ctx, filters, tx); err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif autocommit {\n\t\t\terr = tx.Commit()\n\t\t}\n\t}\n\n\treturn err\n}"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "crud/partials/delete_single.go.tmpl",
		FileModTime: time.Unix(1526939079, 0),
		Content:     string("// Delete deletes a {{.EntityName}} record from database and sets id to nil\nfunc (entity *{{.EntityName}}) Delete(ctx context.Context, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr  error\n\t\tstmt *sql.Stmt\n\t)\n\tid := *entity.ID\n\n\tif tx == nil {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\ttx, err = db.Begin()\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t}\n\t}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tstmt, err = tx.Prepare(\"DELETE FROM {{.Table}} WHERE id = $1\")\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\t{{if .HasPreHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif err = crudPreDelete(id, tx); err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\t_, err = stmt.Exec(id)\n\t\tif err == nil {\n\t\t\tentity.ID = nil\n\t\t} else {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif err = crudPostDelete(id, tx); err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif autocommit {\n\t\t\terr = tx.Commit()\n\t\t}\n\t}\n\t\n\treturn nil\n}"),
	}
	fileb := &embedded.EmbeddedFile{
		Filename:    "crud/partials/get.go.tmpl",
		FileModTime: time.Unix(1526938054, 0),
		Content:     string("// Get returns a single {{.EntityName}} from database by primary key\nfunc Get(ctx context.Context, id {{pkeyPropertyType .PrimaryKey}}) (*{{.EntityName}}, error) {\n\tvar (\n\t\trows   *sql.Rows\n\t\terr    error\n\t\tentity = New()\n\t)\n\t{{if .HasPreHook}}\n    if err = crudPreGet(ctx, id); err != nil {\n\t\treturn nil, err\n\t}\n    {{end}}\n\tselect {\n\tcase <- ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\trows, err = db.Query(\"SELECT {{.SQLFields}} FROM {{.Table}} WHERE id = $1 ORDER BY .id ASC\", id)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\tdefer rows.Close()\n\tif rows.Next() {\n\t\tselect {\n\t\tcase <- ctx.Done():\n\t\t\treturn nil, ctx.Err()\n\t\tdefault:\n\t\t\terr = rows.Scan({{.StructFields}})\n\t\t\tif err != nil {\n\t\t\t\treturn nil, err\n\t\t\t}\n\t\t} \n\t}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tif err = crudPostGet(ctx, entity); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\t{{end}}\n\n\treturn entity, nil\n}"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "crud/partials/insert.go.tmpl",
		FileModTime: time.Unix(1526939198, 0),
		Content:     string("// Insert performs an SQL insert for {{.EntityName}} record and update instance with inserted id.\nfunc (entity *{{.EntityName}}) Insert(ctx context.Context, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\t{{- if pkeyIsAuto .PrimaryKey -}}\n\t\tid  {{pkeyPropertyType .PrimaryKey}}\n\t\t{{- end}}\n\t\terr  error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\ttx, err = db.Begin()\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t}\n\t}\n\t{{range .Before}}{{.}}{{end}}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\t{{if eq .PrimaryKey \"serial\" -}}\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}}) RETURNING id\")\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{else}}\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}})\")\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{- end}}\n\n\t{{if .HasPreHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif err = crudPreSave(\"INSERT\", entity, tx); err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\t{{if eq .PrimaryKey \"serial\" -}}\n\terr = stmt.QueryRow({{.StructFields}}).Scan(&id)\n\tif err == nil {\n\t\tentity.ID = &id\n\t} else {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{else}}\n\t{{if eq .PrimaryKey \"uuid\" -}}\n\tidUUID, err := uuid.NewV4()\n\t\n\tif err == nil {\n\t\tid = idUUID.String()\n\t} else {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t*entity.ID = id\n\t{{- end}}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\t{{range .Relationships}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tstmtRel, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.ThisID}}, {{.ThatID}}) VALUES ($1, $2)\")\n\t\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, relatedID := range *entity.{{.PropertyName}} {\n\t\t_, err = stmtRel.Exec(entity.ID, relatedID)\n\t\tif err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <- ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\tif err := crudPostSave(\"INSERT\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\n\tselect {\n\tcase <- ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif autocommit {\n\t\t\terr = tx.Commit()\n\t\t}\n\t}\n\n\treturn nil\n}"),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "crud/partials/list.go.tmpl",
		FileModTime: time.Unix(1526938054, 0),
		Content:     string("// List returns a slice containing {{.EntityName}} records\nfunc List(ctx context.Context, filters []models.ListFilter) ([]*{{.EntityName}}, error) {\n\tvar (\n\t\tlist     []*{{.EntityName}}\n\t\tsegments []string\n\t\tvalues   []interface{}\n\t\terr      error\n\t\trows     *sql.Rows\n\t)\n\n\tquery := \"SELECT {{.SQLFields}} FROM {{.Table}}\"\n\t{{if .HasPreHook}}\n    if filters, err = crudPreList(ctx, filters); err != nil {\n\t\treturn nil, err\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tfor i, filter := range filters {\n\t\tsegments = append(segments, filter.Field+\" \"+filter.Operation+\" $\"+strconv.Itoa(i+1))\n\t\tvalues = append(values, filter.Value)\n\t}\n\n\tif len(segments) != 0 {\n\t\tquery += \" WHERE \" + strings.Join(segments, \" AND \")\n\t}\n\n\trows, err = db.Query(query+\" ORDER BY id ASC\", values...)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\tdefer rows.Close()\n\tfor rows.Next() {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn nil, ctx.Err()\n\t\tdefault:\n\t\t\tbreak\n\t\t}\n\n\t\tentity := New()\n\t\terr = rows.Scan({{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\tlist = append(list, entity)\n\t}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tif list, err = crudPostList(ctx, list); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\t{{end}}\n\treturn list, nil\n}"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "crud/partials/loadrelated.go.tmpl",
		FileModTime: time.Unix(1526938054, 0),
		Content:     string("// LoadRelated{{.PropertyName}} is a helper function to load related {{.PropertyName}} entities\nfunc LoadRelated{{.PropertyName}}(ctx context.Context, entities []{{.EntityName}}) ([]{{.EntityName}}, error) {\n\tvar (\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices map[{{pkeyPropertyType .PrimaryKey}}]{{trimPrefix .PropertyType \"[]\"}}\n\t\tidthis  {{pkeyPropertyType .PrimaryKey}}\n\t\tidthat  {{trimPrefix .PropertyType \"[]\"}}\n\t\titems   []{{.EntityName}}\n\t)\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tc := 1\n\tfor i, entity := range entities {\n\t\tplaceholder += \"$\" + strconv.Itoa(c) + \",\"\n\t\tindices[entity.ID] = entity\n\t\tvalues = append(values, i)\n\t\tc++\n\t}\n\tplaceholder = strings.TrimRight(placeholder, \",\")\n\n\trows, err := db.Query(\"SELECT {{.ThisID}}, {{.ThatID}} FROM {{.Table}} WHERE {{.ThisID}} IN (\"+placeholder+\")\", values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn nil, ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tfor rows.Next() {\n\t\terr = rows.Scan(&idthis, &idthat)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tindices[idthis].{{.PropertyName}} = append(*indexID[idthis].{{.PropertyName}}, idthat)\n\t}\n\n\treturn nil\n}"),
	}
	filef := &embedded.EmbeddedFile{
		Filename:    "crud/partials/merge.go.tmpl",
		FileModTime: time.Unix(1526938845, 0),
		Content:     string("// Merge performs an SQL merge for {{.Entity.Name}} record.\nfunc (entity *{{.Entity.Name}}) Merge(ctx context.Context, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\ttx, err = db.Begin()\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t}\n\t}\n\n\tif entity.ID == nil {\n\t\treturn entity.Insert(ctx, tx, autocommit)\n\t}\n\n\t{{range .Before}}{{.}}{{end}}\n\n\tstmt, err = tx.Prepare(`INSERT INTO {{.Table}} ({{.SQLFieldsInsert}}) VALUES ({{.SQLPlaceholders}}) \n\tON CONFLICT (id) DO UPDATE SET {{.SQLFieldsUpdate}}`)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\t{{if .HasPreHook}}\n    if err = crudPreSave(\"MERGE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\t_, err = stmt.Exec({{.StructFieldsMerge}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\tif err = crudPostSave(\"MERGE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif autocommit {\n\t\t\terr = tx.Commit()\n\t\t}\n\t}\n\n\treturn err\n}"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "crud/partials/save.go.tmpl",
		FileModTime: time.Unix(1526938054, 0),
		Content:     string("// Save either inserts or updates a {{.EntityName}} record based on whether or not id is nil\nfunc (entity *{{.EntityName}}) Save(ctx context.Context, tx *sql.Tx, autocommit bool) error {\n\t{{if pkeyIsAuto .PrimaryKey -}}\n\tif entity.ID == nil {\n\t\treturn entity.Insert(tx, autocommit)\n\t}\n\treturn entity.Update(tx, autocommit)\n\t{{- else -}}\n\tif entity.ID == nil {\n\t\treturn errors.New(\"primary key cannot be nil\")\n\t}\n\treturn entity.Merge(tx, autocommit)\n\t{{end -}}\n}"),
	}
	fileh := &embedded.EmbeddedFile{
		Filename:    "crud/partials/saverelated.go.tmpl",
		FileModTime: time.Unix(1526938054, 0),
		Content:     string("// SaveRelated{{.PropertyName}} is a helper function to save related {{.PropertyName}} entities\nfunc SaveRelated{{.PropertyName}}(ctx context.Context, idthis  {{pkeyPropertyType .PrimaryKey}}, relatedids {{.PropertyType}}) error {\n\tvar (\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices map[{{pkeyPropertyType .PrimaryKey}}]{{trimPrefix .PropertyType \"[]\"}}\n\t\tidthis  {{pkeyPropertyType .PrimaryKey}}\n\t\tidthat  {{trimPrefix .PropertyType \"[]\"}}\n\t\titems   []{{.EntityName}}\n\t\tstmt    *sql.Stmt\n\t)\n\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\t\n\tstmt, err = tx.Prepare(\"DELETE FROM {{.Table}} WHERE {{.ThisID}} = $1\")\n\n\tif err != nil {\n\t\treturn err\n\t}\n\n\t_, err = stmt.Exec(*entity.ID)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.ThisID}}, {{.ThatID}}) VALUES ($1, $2)\")\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, relatedID := range *entity.{{.PropertyName}} {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\ttx.Rollback()\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\tbreak\n\t\t}\n\n\t\t_, err = stmt.Exec(*entity.ID, relatedID)\n\t\tif err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n}"),
	}
	filei := &embedded.EmbeddedFile{
		Filename:    "crud/partials/update.go.tmpl",
		FileModTime: time.Unix(1526939140, 0),
		Content:     string("// Update Will execute an SQLUpdate Statement for {{.EntityName}} in the database. Prefer using Save instead of Update directly.\nfunc (entity *{{.EntityName}}) Update(ctx context.Context, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\ttx, err = db.Begin()\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t}\n\t}\n\t{{range .Before}}{{.}}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tstmt, err = tx.Prepare(\"UPDATE {{.Table}} SET {{.SQLFields}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn err\n\t}\n\n\t{{if .HasPreHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n    if err = crudPreSave(\"UPDATE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n        return err\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{range .Relationships}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\tstmtRel, err = tx.Prepare(\"DELETE FROM {{.Table}} WHERE {{.ThisID}} = $1\")\n\n\tif err != nil {\n\t\treturn err\n\t}\n\n\t_, err = stmtRel.Exec(*entity.ID)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tstmtRel, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.ThisID}}, {{.ThatID}}) VALUES ($1, $2)\")\n\t\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, relatedID := range *entity.{{.PropertyName}} {\n\t\tselect {\n\t\tcase <-ctx.Done():\n\t\t\ttx.Rollback()\n\t\t\treturn ctx.Err()\n\t\tdefault:\n\t\t\tbreak\n\t\t}\n\n\t\t_, err = stmtRel.Exec(*entity.ID, relatedID)\n\t\tif err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n\t{{end}}\n\t{{if .HasPostHook}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tbreak\n\t}\n\n\tif err = crudPostSave(\"UPDATE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\tselect {\n\tcase <-ctx.Done():\n\t\ttx.Rollback()\n\t\treturn ctx.Err()\n\tdefault:\n\t\tif autocommit {\n\t\t\terr = tx.Commit()\n\t\t}\n\t}\n\n\treturn err\n}"),
	}
	filej := &embedded.EmbeddedFile{
		Filename:    "crud/structure.go.tmpl",
		FileModTime: time.Unix(1526841114, 0),
		Content:     string("package {{.Package}}\n\n// {{.Entity.Name}} {{.Entity.Description}}\ntype {{.Entity.Name}} struct {\n\tID        *{{pkeyPropertyType .Entity.PrimaryKey}}  `json:\"id\"`\n\t{{- range .Entity.Fields}}\n\t{{.Property.Name}} *{{.Property.Type}} `json:\"{{.Serialized}}\"`\n\t{{- end}}\n}\n\n// New returns an instance of {{.Entity.Name}}\nfunc New() *{{.Entity.Name}} {\n\treturn &{{.Entity.Name}}{\n\t\tID: new({{pkeyPropertyType .Entity.PrimaryKey}}),\n\t\t{{- range .Entity.Fields}}\n\t\t{{.Property.Name}}: new({{.Property.Type}}),\n\t\t{{- end}}\n\t}\n}"),
	}
	filek := &embedded.EmbeddedFile{
		Filename:    "http.go.tmpl",
		FileModTime: time.Unix(1526036299, 0),
		Content:     string("package app\n\nimport (\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\n\t\"github.com/gorilla/mux\"\n)\n\n// ServeHTTP starts an http server\nfunc ServeHTTP(listen string, route func(router *mux.Router) error) {\n\tvar err error\n\tsigs := make(chan os.Signal, 1)\n\tsignal.Notify(sigs, syscall.SIGTERM)\n\n\trouter := mux.NewRouter()\n\terr = route(router)\n\n\tif err != nil {\n\t\tlog.Fatal(\"Failed to register routes: \", err)\n\t}\n\n\tgo func() {\n\t\terr = http.ListenAndServe(listen, router)\n\t\tif err != nil {\n\t\t\tlog.Fatal(\"Failed to start http server: \", err)\n\t\t}\n\t}()\n\n\tlog.Println(\"Listening on \" + listen)\n\t<-sigs\n\tlog.Println(\"Server stopped\")\n}\n"),
	}
	filel := &embedded.EmbeddedFile{
		Filename:    "rest.go.tmpl",
		FileModTime: time.Unix(1526781232, 0),
		Content:     string("package {{.Package}}\nimport (\n\t\"database/sql\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"net/http\"\n\t\"strconv\"\n\n\t\"github.com/gorilla/mux\"\n)\n\ntype responseSingle struct {\n\tStatus   bool      `json:\"status\"`\n\tMessages []message `json:\"messages\"`\n\tEntity   *{{.Entity.Name}} `json:\"entity\"`\n}\n\ntype responseList struct {\n\tStatus   bool                `json:\"status\"`\n\tMessages []message           `json:\"messages\"`\n\tEntities []*{{.Entity.Name}} `json:\"entities\"`\n}\n\ntype message struct {\n\tType    rune   `json:\"type\"`\n\tMessage string `json:\"message\"`\n}\n\n//RegisterRoutes registers routes with a mux Router\nfunc RegisterRoutes(router *mux.Router) {\n\t{{if .Entity.Rest.Read}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestGet).Methods(\"GET\"){{end}}\n\t{{if .Entity.Rest.ReadList}}router.HandleFunc(\"/{{.Endpoint}}\", RestList).Methods(\"GET\"){{end}}\n\t{{if .Entity.Rest.Create}}router.HandleFunc(\"/{{.Endpoint}}\", RestCreate).Methods(\"POST\"){{end}}\n\t{{if .Entity.Rest.Update}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestUpdate).Methods(\"PUT\"){{end}}\n\t{{if .Entity.Rest.Delete}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestDelete).Methods(\"DELETE\"){{end}}\n}\n\n{{if .Entity.Rest.Read}}\n//RestGet is a REST endpoint for GET /{{.Endpoint}}/{id}\nfunc RestGet(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\terr      error\n\t\tresponse responseSingle\n\t\t{{if or .Entity.Rest.Hooks.PreRead .Entity.Rest.Hooks.PostRead -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreRead}}\n    if stop, err = restPreGet(w, r, id); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostRead}}\n    if stop, err = restPostGet(w, r, response.Entity); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Status = true\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.ReadList}}\n//RestList is a REST endpoint for GET /{{.Endpoint}}\nfunc RestList(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\tresponse responseList\n\t\tfilters  []models.ListFilter\n\t\t{{if or .Entity.Rest.Hooks.PreList .Entity.Rest.Hooks.PostList}}stop     bool{{end}}\n\t)\n\t{{range .Entity.Fields}}{{if .Filterable}}\n\t{{if eq .Property.Type \"bool\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); val != \"\" {\n\t\tif t, e := strconv.ParseBool(val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:t})\n\t\t}\n\t}\n\t{{end}}\n\t{{if eq .Property.Type \"string\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); val != \"\" {\n\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:val})\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-lk\"); val != \"\" {\n\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"LIKE\", Value:\"%\" + val + \"%\"})\n\t}\n\t{{end}}\n\t{{if eq .Property.Type \"time.Time\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-gt\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\">\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-ge\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\">=\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-lt\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"<\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-le\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"<=\", Value:t})\n\t\t}\n\t}\n\t{{end}}\n\t{{end}}{{end}}\n\n\t{{if .Entity.Rest.Hooks.PreList}}\n    if filters, stop, err = restPreList(w, r, filters); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Entities, err = List(filters)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostList}}\n    if response.Entities, stop, err = restPostList(w, r, response.Entities); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Status = true\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Create}}\n//RestCreate is a REST endpoint for POST /{{.Endpoint}}\nfunc RestCreate(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\trawbody  []byte\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreCreate .Entity.Rest.Hooks.PostCreate}}stop     bool{{end}}\n\t)\n\n\trawbody, err = ioutil.ReadAll(r.Body)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to read body\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity = New()\n\terr = json.Unmarshal(rawbody, response.Entity)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to decode body\"}]}`)\n\t\treturn\n\t}\n\t{{if pkeyIsAuto .Entity.PrimaryKey -}}\n\tresponse.Entity.ID = nil\n\t{{- end}}\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreCreate}}\n\tif stop, err = restPreCreate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\n\terr = response.Entity.Save(tx, false)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Save failed\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostCreate}}\n\tif stop, err = restPostCreate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\t\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestCreate could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Update}}\n//RestUpdate is a REST endpoint for PUT /{{.Endpoint}}/{id}\nfunc RestUpdate(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\trawbody  []byte\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreUpdate .Entity.Rest.Hooks.PostUpdate -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\trawbody, err = ioutil.ReadAll(r.Body)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to read body\"}]}`)\n\t\treturn\n\t}\n\n\terr = json.Unmarshal(rawbody, response.Entity)\n\tif err != nil {\n\t\tif err != nil {\n\t\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\t\tw.WriteHeader(http.StatusBadRequest)\n\t\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to decode body\"}]}`)\n\t\t\treturn\n\t\t}\n\t}\n\tresponse.Entity.ID = &id\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreUpdate}}\n    if stop, err = restPreUpdate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n        return\n    } else if stop {\n\t\treturn\n\t}\n    {{end}}\n\n\terr = response.Entity.Save(tx, false)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Save failed\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostUpdate}}\n    if stop, err = restPostUpdate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n        return\n    } else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\t\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestUpdate could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Delete}}\n//RestDelete is a REST endpoint for DELETE /{{.Endpoint}}/{id}\nfunc RestDelete(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\terr      error\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreDelete .Entity.Rest.Hooks.PostDelete -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\t{{if .Entity.Rest.Hooks.PreDelete}}\n\tif stop, err = restPreDelete(w, r, id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n    {{end}}\n\terr = response.Entity.Delete(tx, false)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Delete failed\"}]}`)\n\t\treturn\n\t}\n\t{{if .Entity.Rest.Hooks.PostDelete}}\n\tif stop, err = restPostDelete(w, r, id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestDelete could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n"),
	}
	filem := &embedded.EmbeddedFile{
		Filename:    "rest_hooks.go.tmpl",
		FileModTime: time.Unix(1526025578, 0),
		Content:     string("package {{.Package}}\n\nimport (\n\t\"database/sql\"\n\t\"net/http\"\n)\n\n{{if .Hooks.PreRead}}\nfunc restPreGet(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostRead}}\nfunc restPostGet(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreList}}\nfunc restPreList(w http.ResponseWriter, r *http.Request, filters []models.ListFilter) ([]models.ListFilter, bool, error) {\n\treturn filters, false, nil\n}\n{{end}}\n{{if .Hooks.PostList}}\nfunc restPostList(w http.ResponseWriter, r *http.Request, list []*{{.Entity.Name}}) ([]*{{.Entity.Name}}, bool, error) {\n\treturn list, false, nil\n}\n{{end}}\n\n{{if .Hooks.PreCreate}}\nfunc restPreCreate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostCreate}}\nfunc restPostCreate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreUpdate}}\nfunc restPreUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostUpdate}}\nfunc restPostUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreDelete}}\nfunc restPreDelete(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostDelete}}\nfunc restPostDelete(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}"),
	}
	filen := &embedded.EmbeddedFile{
		Filename:    "schema.sql.tmpl",
		FileModTime: time.Unix(1526359511, 0),
		Content:     string("DROP TABLE IF EXISTS {{.Entity.Table}};\n\nCREATE TABLE {{.Entity.Table}} (\n\t\"id\" {{pkeyFieldType .Entity.PrimaryKey}},\n\t{{- range $i, $e := .Entity.Fields}}{{if ne .Schema.Field \"\"}}\n\t\"{{.Schema.Field}}\" {{$e.Schema.Type}}\n\t{{- if not .Schema.Nullable}} NOT NULL{{end}}\n\t{{- if ne .Schema.Default \"\" -}} DEFAULT {{.Schema.Default}}{{end}},\n\t{{- end}}{{- end}}\n\t{{range .Entity.TableConstraints}}{{.}},{{end}}\n\tPRIMARY KEY (\"id\")\n);\n\n{{range .ManyManyFields}}\nDROP TABLE IF EXISTS {{.Relationship.Target.Table}};\n\nCREATE TABLE {{.Relationship.Target.Table}} (\n\t\"{{.Relationship.Target.ThisID}}\" {{pkeyFieldType $.Entity.PrimaryKey}} NOT NULL,\n\t\"{{.Relationship.Target.ThatID}}\" {{.Relationship.ThatFieldType}} NOT NULL\n);\n{{end}}"),
	}
	fileo := &embedded.EmbeddedFile{
		Filename:    "vuetify_actions.js.tmpl",
		FileModTime: time.Unix(1525883267, 0),
		Content:     string("import types from \"./types\";\n\nexport default {}"),
	}
	filep := &embedded.EmbeddedFile{
		Filename:    "vuetify_edit.vue.tmpl",
		FileModTime: time.Unix(1525372377, 0),
		Content:     string("<template>\n    <div class=\"container\">\n\t\t<v-toolbar color=\"transparent\" flat>\n            <v-toolbar-title class=\"grey--text text--darken-4 ml-0\"><h2>{{.Entity.Name}}</h2></v-toolbar-title>\n            <v-spacer></v-spacer>\n            <v-btn ml-0 small color=\"grey\" flat :to=\"{name: '{{.Endpoint}}List'}\">\n                <v-icon dark>arrow_back</v-icon> Back\n            </v-btn>\n        </v-toolbar>\n\t\t<v-alert :type=\"message.type\" :value=\"true\" v-for=\"(message, index) in messages\" :key=\"index\">\n\t\t{{ \"{{ message.text }}\" }}\n\t\t</v-alert>\n  \n        {{range .Entity.Fields -}}\n        {{widget_field \"vuetify\" .Widget.Type .}}\n        {{- end}}\n\n        <v-btn color=\"primary\" @click=\"save()\">Save</v-btn>\n        <v-btn color=\"gray\" :to=\"{name: '{{.Endpoint}}List'}\">Cancel</v-btn>\n\t</div>\n</template>\n  \n<script>\nimport axios from \"axios\"\n\nexport default {\n    props: [\"id\"],\n    created() {\n        if (!this.id) {\n            return\n        }\n\n        axios.get(\"/api/{{.Endpoint}}/\" + this.id).then(response => {\n            this.id = response.data.entity.id\n            {{range .Entity.Fields}}{{if ne .Serialized \"id\"}}\n            this.entity.{{.Serialized}} = response.data.entity.{{.Serialized}}\n            {{if eq .Widget.Type \"date\"}}this.dates.{{.Serialized}}.value = response.data.entity.{{.Serialized}}.substr(0,10){{end}}\n            {{end}}{{end}}\n        })\n    },\n    data() {\n        return {\n            select: {\n                {{range $i, $v := .Entity.Fields}}{{if eq .Widget.Type \"date\"}}\n                {{.Serialized}}: {\n                    items:[\n                        {{range $j, $u := .Widget.Options}}\n                        {text: \"{{.Label}}\", value: \"{{.Value}}\"}{{if eq (plus1 $j) (len $u)}},{{end}}\n                        {{end}}\n                    ]\n                }{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n\t\t\t},\n\t\t\tdates: {\n                {{range $i, $v := .Entity.Fields}}{{if eq .Widget.Type \"date\"}}\n\t\t\t\t{{.Serialized}}: {value: null, menuAppear: false}{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n\t\t\t},\n            messages: [],\n            entity: {\n                {{range $i, $e := .Entity.Fields}}{{if ne .Serialized \"id\"}}\n                {{.Serialized}} : null{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n            }\n        }\n    },\n    watch: {\n        {{range $i, $e := .Entity.Fields}}\n        \"select.{{.Serialized}}.search\": function(val) {\n            val && this.querySelections(\"{{.Serialized}}\", \"{{$.Endpoint}}\", \"{{$.Prefix}}{{.Relationship.Target.Endpoint}}\", val)\n        }{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n        {{end}}\n    },\n    methods: {\n        querySelections(fieldname, endpoint, filter, val) {\n            this.select[fieldname].loading = true\n            axios.get(\"/api/\" + endpoint + \"?\" + filter + \"-lk=\" + encodeURIComponent(val)).then(response => {\n                this.select[fieldname].loading = false\n                this.select[fieldname].items = response.data.entities.map(function(e) {\n                    return { text: e[filter], value: e.id }\n                })\n            })\n        },\n        save() {\n            if (this.id) {\n                axios.put(\"/api/{{.Endpoint}}/\" + this.id, this.entity).then(this.saved)\n            } else {\n                axios.post(\"/api/{{.Endpoint}}\", this.entity).then(this.saved)\n            }\n\t\t},\n\t\tsaved(response) {\n\t\t\tthis.id = response.data.entity.id\n\t\t\t{{range .Entity.Fields}}{{if ne .Serialized \"id\"}}\n            this.entity.{{.Serialized}} = response.data.entity.{{.Serialized}}\n            {{end}}{{end}}\n\n\t\t\tthis.messages.push({\n\t\t\t\ttype: \"success\",\n\t\t\t\ttext: \"{{.Entity.Name}} saved successfully\"\n\t\t\t})\n\t\t}\n    }\n}\n</script>"),
	}
	fileq := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-checkbox.vue.tmpl",
		FileModTime: time.Unix(1524654324, 0),
		Content:     string("<v-checkbox label=\"{{.Label}}\" v-model=\"entity.{{.Serialized}}\"></v-checkbox>"),
	}
	filer := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-date.vue.tmpl",
		FileModTime: time.Unix(1524654335, 0),
		Content:     string("<v-menu\n\tref=\"menu_{{.Serialized}}\"\n\tlazy\n\t:close-on-content-click=\"false\"\n\tv-model=\"dates.{{.Serialized}}.menuAppear\"\n\ttransition=\"scale-transition\"\n\toffset-y\n\tfull-width\n\t:nudge-right=\"40\"\n\tmin-width=\"290px\"\n\t:return-value.sync=\"dates.{{.Serialized}}.value\"\n\t>\n\t<v-text-field\n\t\tslot=\"activator\"\n\t\tlabel=\"{{.Label}}\"\n\t\tv-model=\"dates.{{.Serialized}}.value\"\n\t\tprepend-icon=\"event\"\n\t\treadonly\n\t\t></v-text-field>\n\t\t<v-date-picker v-model=\"dates.{{.Serialized}}.value\" @change=\"entity.{{.Serialized}} = dates.{{.Serialized}}.value + 'T00:00:00Z'\" no-title scrollable>\n\t\t<v-spacer></v-spacer>\n\t\t<v-btn flat color=\"primary\" @click=\"menu_{{.Serialized}} = false\">Cancel</v-btn>\n\t\t<v-btn flat color=\"primary\" @click=\"$refs.menu_{{.Serialized}}.save(dates.{{.Serialized}}.value)\">OK</v-btn>\n\t\t</v-date-picker>\n</v-menu>"),
	}
	files := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-number.vue.tmpl",
		FileModTime: time.Unix(1524654354, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" type=\"number\" />"),
	}
	filet := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-password.vue.tmpl",
		FileModTime: time.Unix(1524654359, 0),
		Content:     string("<v-text-field\n\tv-model=\"entity.{{.Serialized}}\"\n\t:append-icon=\"e1 ? 'visibility' : 'visibility_off'\"\n\t:append-icon-cb=\"() => (e1 = !e1)\"\n\t:type=\"e1 ? 'password' : 'text'\"\n\tcounter\n  ></v-text-field>"),
	}
	fileu := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-select-rel.vue.tmpl",
		FileModTime: time.Unix(1524663870, 0),
		Content:     string("<v-select\n    autocomplete\n    cache-items\n    required\n    label=\"{{.Label}}\"\n    :loading=\"select.{{.Serialized}}.isloading\"\n    :items=\"select.{{.Serialized}}.items\"\n\t:search-input.sync=\"select.{{.Serialized}}.search\"\n\tv-model=\"entity.{{.Serialized}}\"\n\t{{if .Widget.Multiple}}multiple chips{{end}}\n></v-select>"),
	}
	filev := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-select.vue.tmpl",
		FileModTime: time.Unix(1524663879, 0),
		Content:     string("<v-select\n\tautocomplete\n\tcache-items\n\trequired\n\tlabel=\"{{.Label}}\"\n\t:items=\"select.{{.Serialized}}.items\"\n\tv-model=\"entity.{{.Serialized}}\"\n\t{{if .Widget.Multiple}}multiple chips{{end}}\n></v-select>"),
	}
	filew := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-textarea.vue.tmpl",
		FileModTime: time.Unix(1524654397, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" multiline />"),
	}
	filex := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-textfield.vue.tmpl",
		FileModTime: time.Unix(1524654403, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" />"),
	}
	filey := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-time.vue.tmpl",
		FileModTime: time.Unix(1524654414, 0),
		Content:     string("<div>\n\t<v-time-picker v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" :landscape=\"landscape\"></v-time-picker>\n</div>"),
	}
	filez := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-toggle.vue.tmpl",
		FileModTime: time.Unix(1524654423, 0),
		Content:     string("<v-switch\n\tlabel=\"{{.Label}}\"\n\tv-model=\"entity.{{.Serialized}}\"\n></v-switch>"),
	}
	file10 := &embedded.EmbeddedFile{
		Filename:    "vuetify_getters.js.tmpl",
		FileModTime: time.Unix(1524749464, 0),
		Content:     string("export default {}\n"),
	}
	file11 := &embedded.EmbeddedFile{
		Filename:    "vuetify_index.js.tmpl",
		FileModTime: time.Unix(1524774338, 0),
		Content:     string("import actions from \"./actions\";\nimport getters from \"./getters\";\nimport mutations from \"./mutations\";\nimport routes from \"./routes\";\n\nconst namespaced = true;\n\nconst state = {\n  entities: routes.routes\n};\n\nexport default {\n  namespaced,\n  state,\n  actions,\n  getters,\n  mutations\n};"),
	}
	file12 := &embedded.EmbeddedFile{
		Filename:    "vuetify_list.vue.tmpl",
		FileModTime: time.Unix(1524663450, 0),
		Content:     string("<template>\n    <v-container>\n        <v-toolbar color=\"transparent\" flat>\n            <v-toolbar-title class=\"grey--text text--darken-4 ml-0\"><h2>{{.Entity.Name}}</h2></v-toolbar-title>\n            <v-spacer></v-spacer>\n            <v-btn mr-0 color=\"primary\" :to=\"{name: '{{.Endpoint}}Edit', params:{id: 0}}\">\n                <v-icon dark>add</v-icon> Add\n            </v-btn>\n        </v-toolbar>\n        \n        <v-alert :type=\"message.type === 'E' ? 'error' : message.type\" :value=\"true\" v-for=\"(message, index) in messages\" :key=\"index\">\n            {{ \"{{ message.text }}\" }}\n        </v-alert>\n\n        <v-alert type=\"info\" value=\"true\"  color=\"primary\" outline icon=\"info\" v-if=\"entities.length === 0\">\n            No {{.Entity.Name}} exist. Would you like to create one now?\n            <v-btn :to=\"{name: '{{.Endpoint}}Edit', params:{id: 0}}\" color=\"primary\">create new</v-btn>\n        </v-alert>\n        <template v-else>\n            <v-text-field mb-4 append-icon=\"search\" label=\"Search\" single-line hide-details v-model=\"search\"></v-text-field>            \n            <v-data-table :headers=\"headers\" :items=\"entities\" class=\"elevation-1\" :search=\"search\">\n                <template slot=\"items\" slot-scope=\"props\">\n\t\t\t\t\t{{ range .Entity.Fields }}\n\t\t\t\t\t<td>{{ printf \"{{ props.item.%s}}\" .Serialized }}</td>\n\t\t\t\t\t{{end}}\n                    <td class=\"justify-center layout px-0\">\n                        <v-btn icon class=\"mx-0\" :to=\"{name: '{{.Endpoint}}Edit', params: {'id': props.item.id}  }\">\n                            <v-icon color=\"teal\">edit</v-icon>\n                        </v-btn>\n                    </td>\n                </template>\n\n                <template slot=\"no-data\">\n                    <v-flex ma-4>\n                        <v-alert slot=\"no-results\" :value=\"true\" color=\"primary\" outline icon=\"info\" v-if=\"search.length > 0\">\n                        Your search for \"{{ \"{{ search }}\" }}\" found no results.\n                        </v-alert>\n                        <v-alert slot=\"no-results\" :value=\"true\" color=\"primary\" outline icon=\"info\" v-else>\n                            No {{.Entity.Name}} found.\n                        </v-alert>\n                    </v-flex>\n                </template>\n            </v-data-table>\n        </template>\n    </v-container>\n</template>\n\n<script>\nimport axios from \"axios\"\nexport default {\n  data() {\n    return {\n      messages: [],\n      search: \"\",\n      headers: [\n\t\t{{range .Entity.Fields }}\n\t\t{text: \"{{.Label}}\", value: \"{{.Serialized}}\"},\n\t\t{{end}}\n        {'text': 'Action', 'value': null}\n      ],\n      entities: []\n    };\n  },\n  created() {\n    axios\n      .get(\"/api/{{.Endpoint}}\")\n      .then(response => {\n        this.entities = response.data.entities;\n      })\n      .catch(error => {\n        this.messages = [...this.messages, ...error.response.data.messages];\n      });\n  }\n};\n</script>"),
	}
	file13 := &embedded.EmbeddedFile{
		Filename:    "vuetify_mutations.js.tmpl",
		FileModTime: time.Unix(1524749492, 0),
		Content:     string("import types from \"./types\";\n\nexport default {}\n"),
	}
	file14 := &embedded.EmbeddedFile{
		Filename:    "vuetify_routes.js.tmpl",
		FileModTime: time.Unix(1524774098, 0),
		Content:     string("{{range .Entities}}\n// {{.Name}} {{.Description}}\nimport {{.Name}}Edit from \"../views/{{plural .Name}}Edit.vue\";\nimport {{.Name}}List from \"../views/{{plural .Name}}List.vue\";\n{{end}}\n\nlet routes = [\n  {{range $i, $v := .Entities}}\n  {\n    path: \"/{{lower (plural .Name)}}/:id\",\n    name: \"{{lower (plural .Name)}}Edit\",\n    props: true,\n    icon: \"dashboard\",\n    component: {{.Name}}Edit,\n    entity: \"{{plural .Name}}\"\n  },\n  {\n    path: \"/{{lower (plural .Name)}}list/\",\n    name: \"{{lower (plural .Name)}}List\",\n    icon: \"dashboard\",\n    component: {{.Name}}List,\n    entity: \"{{plural .Name}}\"\n  }{{if ne (plus1 $i) (len $.Entities)}},{{end}}\n  {{end}}\n];\n\nlet entities = [\n  {{range $i, $v := .Entities}}\n  \"{{plural .Name}}\"{{if ne (plus1 $i) (len $.Entities)}},{{end}}\n  {{end}}\n];\n\nfunction registerRoutes(router) {\n  router.addRoutes(routes);\n}\n\nexport default {\n  routes,\n  entities,\n  registerRoutes\n}\n"),
	}
	file15 := &embedded.EmbeddedFile{
		Filename:    "vuetify_types.js.tmpl",
		FileModTime: time.Unix(1524749554, 0),
		Content:     string("export default {}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1526937997, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2,  // "bootstrap.go.tmpl"
			file3,  // "bootstrap_env.tmpl"
			filek,  // "http.go.tmpl"
			filel,  // "rest.go.tmpl"
			filem,  // "rest_hooks.go.tmpl"
			filen,  // "schema.sql.tmpl"
			fileo,  // "vuetify_actions.js.tmpl"
			filep,  // "vuetify_edit.vue.tmpl"
			fileq,  // "vuetify_editor-field-checkbox.vue.tmpl"
			filer,  // "vuetify_editor-field-date.vue.tmpl"
			files,  // "vuetify_editor-field-number.vue.tmpl"
			filet,  // "vuetify_editor-field-password.vue.tmpl"
			fileu,  // "vuetify_editor-field-select-rel.vue.tmpl"
			filev,  // "vuetify_editor-field-select.vue.tmpl"
			filew,  // "vuetify_editor-field-textarea.vue.tmpl"
			filex,  // "vuetify_editor-field-textfield.vue.tmpl"
			filey,  // "vuetify_editor-field-time.vue.tmpl"
			filez,  // "vuetify_editor-field-toggle.vue.tmpl"
			file10, // "vuetify_getters.js.tmpl"
			file11, // "vuetify_index.js.tmpl"
			file12, // "vuetify_list.vue.tmpl"
			file13, // "vuetify_mutations.js.tmpl"
			file14, // "vuetify_routes.js.tmpl"
			file15, // "vuetify_types.js.tmpl"

		},
	}
	dir4 := &embedded.EmbeddedDir{
		Filename:   "crud",
		DirModTime: time.Unix(1526938273, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file5, // "crud/crud.go.tmpl"
			file6, // "crud/filters.go.tmpl"
			file7, // "crud/hooks.go.tmpl"
			filej, // "crud/structure.go.tmpl"

		},
	}
	dir8 := &embedded.EmbeddedDir{
		Filename:   "crud/partials",
		DirModTime: time.Unix(1526939079, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file9, // "crud/partials/delete_many.go.tmpl"
			filea, // "crud/partials/delete_single.go.tmpl"
			fileb, // "crud/partials/get.go.tmpl"
			filec, // "crud/partials/insert.go.tmpl"
			filed, // "crud/partials/list.go.tmpl"
			filee, // "crud/partials/loadrelated.go.tmpl"
			filef, // "crud/partials/merge.go.tmpl"
			fileg, // "crud/partials/save.go.tmpl"
			fileh, // "crud/partials/saverelated.go.tmpl"
			filei, // "crud/partials/update.go.tmpl"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir4, // "crud"

	}
	dir4.ChildDirs = []*embedded.EmbeddedDir{
		dir8, // "crud/partials"

	}
	dir8.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1526937997, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":              dir1,
			"crud":          dir4,
			"crud/partials": dir8,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"bootstrap.go.tmpl":                        file2,
			"bootstrap_env.tmpl":                       file3,
			"crud/crud.go.tmpl":                        file5,
			"crud/filters.go.tmpl":                     file6,
			"crud/hooks.go.tmpl":                       file7,
			"crud/partials/delete_many.go.tmpl":        file9,
			"crud/partials/delete_single.go.tmpl":      filea,
			"crud/partials/get.go.tmpl":                fileb,
			"crud/partials/insert.go.tmpl":             filec,
			"crud/partials/list.go.tmpl":               filed,
			"crud/partials/loadrelated.go.tmpl":        filee,
			"crud/partials/merge.go.tmpl":              filef,
			"crud/partials/save.go.tmpl":               fileg,
			"crud/partials/saverelated.go.tmpl":        fileh,
			"crud/partials/update.go.tmpl":             filei,
			"crud/structure.go.tmpl":                   filej,
			"http.go.tmpl":                             filek,
			"rest.go.tmpl":                             filel,
			"rest_hooks.go.tmpl":                       filem,
			"schema.sql.tmpl":                          filen,
			"vuetify_actions.js.tmpl":                  fileo,
			"vuetify_edit.vue.tmpl":                    filep,
			"vuetify_editor-field-checkbox.vue.tmpl":   fileq,
			"vuetify_editor-field-date.vue.tmpl":       filer,
			"vuetify_editor-field-number.vue.tmpl":     files,
			"vuetify_editor-field-password.vue.tmpl":   filet,
			"vuetify_editor-field-select-rel.vue.tmpl": fileu,
			"vuetify_editor-field-select.vue.tmpl":     filev,
			"vuetify_editor-field-textarea.vue.tmpl":   filew,
			"vuetify_editor-field-textfield.vue.tmpl":  filex,
			"vuetify_editor-field-time.vue.tmpl":       filey,
			"vuetify_editor-field-toggle.vue.tmpl":     filez,
			"vuetify_getters.js.tmpl":                  file10,
			"vuetify_index.js.tmpl":                    file11,
			"vuetify_list.vue.tmpl":                    file12,
			"vuetify_mutations.js.tmpl":                file13,
			"vuetify_routes.js.tmpl":                   file14,
			"vuetify_types.js.tmpl":                    file15,
		},
	})
}
