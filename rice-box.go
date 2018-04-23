package main

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "bootstrap.go.tmpl",
		FileModTime: time.Unix(1524249429, 0),
		Content:     string("package main\n\nimport (\n\t\"database/sql\"\n\t\"log\"\n\t\"os\"\n\n\t\"github.com/joho/godotenv\"\n\t_ \"github.com/lib/pq\"\n)\n\nconst (\n\t//EnvironmentProd represents production environment\n\tEnvironmentProd = \"PROD\"\n\n\t//EnvironmentDev represents development environment\n\tEnvironmentDev  = \"DEV\"\n)\n\nvar (\n\tenv string\n    db  *sql.DB\n\t{{range .Settings}}\n\t// {{.Name}} {{.Description}}\n\t{{.Name}} {{.Type}}\n\t{{end}}\n)\n\nfunc bootstrap() {\n\tvar err error\n\n\tgodotenv.Load()\n\n\tdsn := os.Getenv(\"DSN\")\n\tenv = os.Getenv(\"ENV\")\n\n\tif env == \"\" {\n\t\tenv = EnvironmentProd\n\t}\n\n\tif dsn == \"\" {\n\t\tlog.Fatal(\"Environment variable DSN must be defined. Example: postgres://user:pass@host/db?sslmode=disable\")\n\t}\n\n\tdb, err = sql.Open(\"postgres\", dsn)\n\tif err == nil {\n\t\tlog.Println(\"Connected to database successfully.\")\n\t} else if (env == EnvironmentDev) {\n\t\tlog.Println(\"Database connection failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database connection failed: \", err)\n\t}\n\n\terr = db.Ping()\n\tif err == nil {\n\t\tlog.Println(\"Pinged database successfully.\")\n\t} else if (env == EnvironmentDev) {\n\t\tlog.Println(\"Database ping failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database ping failed: \", err)\n\t}\n}"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "crud.go.tmpl",
		FileModTime: time.Unix(1524466573, 0),
		Content:     string("package x\n\nimport \"database/sql\"\n\nvar db *sql.DB\n\n// Inject allows injection of services into the package\nfunc Inject(database *sql.DB) {\n\tdb = database\n}\n\n//New return a new {{.Name}} instance\nfunc New() *{{.Name}} {\n\tentity := new({{.Name}})\n\t{{.Fields}}\n\n\treturn entity\n}\n\n// List returns a slice containing {{.Name}} records\nfunc List(filters []models.ListFilter) ([]*{{.Name}}, error) {\n\tvar (\n\t\tlist     []*{{.Name}}\n\t\tsegments []string\n\t\tvalues   []interface{}\n\t\terr      error\n\t)\n\n\tquery := \"SELECT {{.SQLFields}} FROM {{.TableName}}\"\n\t{{if .PreExecHook }}\n    if filters, err = crudPreList(filters); err != nil {\n\t\treturn nil, fmt.Errorf(\"error executing crudPreList() in List(filters) for entity '{{.Name}}': %s\", err)\n\t}\n    {{end}}\n\tfor i, filter := range filters {\n\t\tsegments = append(segments, filter.Field+\" \"+filter.Operation+\" $\"+strconv.Itoa(i+1))\n\t\tvalues = append(values, filter.Value)\n\t}\n\n\tif len(segments) != 0 {\n\t\tquery += \" WHERE \" + strings.Join(segments, \" AND \")\n\t}\n\n\trows, err := db.Query(query+\" ORDER BY id ASC\", values...) {{.ManyIndexDecl}}\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\tdefer rows.Close()\n\tfor rows.Next() {\n\t\tentity := New()\n\t\terr := rows.Scan({{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\tlist = append(list, entity) {{.ManyIndexAssign}}\n\t}\n\t{{.ManyMany}}\n\t{{if .PostExecHook }}\n\tif list, err = crudPostList(list); err != nil {\n\t\treturn nil, fmt.Errorf(\"error executing crudPostList() in List(filters) for entity '{{.Name}}': %s\", err)\n\t}\n\t{{end}}\n\treturn list, nil\n}\n{{.ManyManyLoadRelated}}\n\n// Delete deletes a {{.Name}} record from database by id primary key\nfunc Delete(id int64, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr error\n\t\t{{.ManyManyVars}}\n\t)\n\n\tif tx == nil {\n\t\ttx, err = db.Begin()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\tstmt, err := tx.Prepare(\"DELETE FROM {{.TableName}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{if .PreExecHook}}\n\tif err := crudPreDelete(id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPreDelete() in Delete(%d) for entity '{{.Name}}': %s\", id, err)\n\t}\n\t{{end}}\n\t{{.ManyManyFunc}}\n\t_, err = stmt.Exec(id)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing transaction statement in Delete(%d) for entity '{{.Name}}': %s\", id, err)\n\t}\n\t{{if .PostExecHook}}\n\tif err := crudPostDelete(id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"Error executing crudPostDelete() in Delete(%d) for entity '{{.Name}}': %s\", id, err)\n\t}\n\t{{end}}\n\tif autocommit {\n\t\terr = tx.Commit()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error committing transaction in Delete(%d) for '{{.Name}}': %s\", id, err)\n\t\t}\n\t}\n\n\treturn err\n}\n\n// Delete deletes a {{.Name}} record from database and sets id to nil\nfunc (entity *{{.Name}}) Delete(tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr error\n\t\t{{.ManyManyVars}}\n\t)\n\n\tid := *entity.ID\n\n\tif tx == nil {\n\t\ttx, err = db.Begin()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\tstmt, err := tx.Prepare(\"DELETE FROM {{.TableName}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{if .PreExecHook}}\n\tif err := crudPreDelete(id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPreDelete() in {{.Name}}.Delete() for ID = %d : %s\", id, err)\n\t}\n\t{{end}}\n\t{{.ManyManyMethod}}\n\t_, err = stmt.Exec(id)\n\tif err == nil {\n\t\tentity.ID = nil\n\t} else {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing transaction statement in {{.Name}}.Delete() for ID = %d : %s\", id, err)\n\t}\n\t{{if .PostExecHook}}\n\tif err = crudPostDelete(id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPostDelete() in {{.Name}}.Delete() for ID = %d : %s\", id, err)\n\t}\n\t{{end}}\n\tif autocommit {\n\t\terr = tx.Commit()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error committing transaction in {{.Name}}.Delete() for ID = %d : %s\", id, err)\n\t\t}\n\t}\n\n\treturn err\n}\n\n// Save either inserts or updates a {{.Name}} record based on whether or not id is nil\nfunc (entity *{{.Name}}) Save(tx *sql.Tx, autocommit bool) error {\n\tif entity.ID == nil {\n\t\treturn entity.Insert(tx, autocommit)\n\t}\n\treturn entity.Update(tx, autocommit)\n}\n\n// Insert performs an SQL insert for {{.Name}} record and update instance with inserted id.\n// Prefer using Save rather than Insert directly.\nfunc (entity *{{.Name}}) Insert(tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\tid  int64\n\t\terr error\n\t\t{{.ManyManyVars}}\n\t)\n\n\tif tx == nil {\n\t\ttx, err = db.Begin()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\t{{.SnippetsBefore}}\n\t\n\tstmt, err := tx.Prepare(\"INSERT INTO {{.TableName}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}}) RETURNING id\")\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{if .PreExecHook }}\n    if err := crudPreSave(entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPreSave() in {{.Name}}.Insert(): %s\", err)\n\t}\n    {{end}}\n\terr = stmt.QueryRow({{.StructFields}}).Scan(&id)\n\tif err == nil {\n\t\tentity.ID = &id\n\t} else {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing transaction statement in {{.Name}}: %s\", err)\n\t}\n\t{{.ManyMany}}\n\t{{if .PostExecHook }}\n\tif err := crudPostSave(entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPostSave() in {{.Name}}.Insert(): %s\", err)\n\t}\n\t{{end}}\n\tif autocommit {\n\t\terr = tx.Commit()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error committing transaction in {{.Name}}.Insert(): %s\", err)\n\t\t}\n\t}\n\n\treturn nil\n}\n\n//Update Will execute an SQLUpdate Statement for {{.Name}} in the database. Prefer using Save instead of Update directly.\nfunc (entity *{{.Name}}) Update(tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr error\n\t\t{{.ManyManyVars}}\n\t)\n\n\tif tx == nil {\n\t\ttx, err = db.Begin()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\t{{.SnippetsBefore}}\n\n\tstmt, err := tx.Prepare(\"UPDATE {{.TableName}} SET {{.SQLFields}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn err\n\t}\n\n\t{{if .PreExecHook }}\n    if err := crudPreSave(entity, tx); err != nil {\n\t\ttx.Rollback()\n        return fmt.Errorf(\"error executing crudPreSave() in {{.Name}}.Update(): %s\", err)\n\t}\n    {{end}}\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing transaction statement in {{.Name}}.Update(): %s\", err)\n\t}\n\t{{.ManyMany}}\n\t{{if .PostExecHook }}\n\tif err := crudPostSave(entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn fmt.Errorf(\"error executing crudPostSave() in {{.Name}}.Update(): %s\", err)\n\t}\n\t{{end}}\n\tif autocommit {\n\t\terr = tx.Commit()\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"error committing transaction in {{.Name}}.Update(): %s\", err)\n\t\t}\n\t}\n\n\treturn nil\n}"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "http.go.tmpl",
		FileModTime: time.Unix(1524431103, 0),
		Content:     string("package main\n\nimport (\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\n\t\"github.com/gorilla/mux\"\n)\n\n// serve starts an http server\nfunc serve(route func(prefix string, router *mux.Router) error) {\n\tvar err error\n\tsigs := make(chan os.Signal, 1)\n\tsignal.Notify(sigs, syscall.SIGTERM)\n\n\trouter := mux.NewRouter()\n\terr = route(\"{{.Prefix}}\", router)\n\n\tif err != nil {\n\t\tlog.Fatal(\"Failed to register routes: \", err)\n\t}\n\t\n\tgo func() {\n\t\terr = http.ListenAndServe(\":{{.Port}}\", router)\n\t\tif err != nil {\n\t\t\tlog.Fatal(\"Failed to start http server: \", err)\n\t\t}\n\t}()\n\n\tlog.Println(\"Listening on : {{.Port}}\")\n\t<-sigs\n\tlog.Println(\"Server stopped\")\n}"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1524466098, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "bootstrap.go.tmpl"
			file3, // "crud.go.tmpl"
			file4, // "http.go.tmpl"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1524466098, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"bootstrap.go.tmpl": file2,
			"crud.go.tmpl":      file3,
			"http.go.tmpl":      file4,
		},
	})
}
