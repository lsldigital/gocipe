package main

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "application/gen-service.sh.tmpl",
		FileModTime: time.Unix(1527838716, 0),
		Content:     string("#!/bin/bash\n\nCURRENT_DIR=`pwd`\n\nif [ \"$1\" == \"\" ]; then\n    echo \"Enter a valid filename\"\n\texit 1\nfi\n\nif [ ! -d \"${CURRENT_DIR}/vendor\" ]; then\n    echo \"vendor folder not found in ${CURRENT_DIR}\"\n    echo \"Running dep init\"\n    dep init\nfi\n\nif [ ! -f \"${CURRENT_DIR}/web/package.json\" ]; then\n    echo \"package.json file not found in ${CURRENT_DIR}/web\"\n    echo \"Running npm init\"\n    cd \"${CURRENT_DIR}/web\"\n    npm init --yes\nfi\n\nif [ ! -d \"${CURRENT_DIR}/web/node_modules/\" ]; then\n    echo \"node_modules dir not found in ${CURRENT_DIR}/web\"\n    echo \"Running npm init\"\n    cd \"${CURRENT_DIR}/web\"\n    npm init --yes\nfi\n\nif [ ! -d \"${CURRENT_DIR}/web/node_modules/.bin\" ]; then\n    echo \".bin not found dir in ${CURRENT_DIR}/web/node_modules\"\n    echo \"Running npm init\"\n    cd \"${CURRENT_DIR}/web\"\n    npm init --yes\nfi\n\nif [ ! -f \"${CURRENT_DIR}/web/node_modules/.bin/protoc-gen-ts\" ]; then\n    echo \"protoc-gen-ts executable not found in ${CURRENT_DIR}/web/node_modules/.bin\"\n    echo \"Running npm install --save-dev ts-protoc-gen\"\n    cd \"${CURRENT_DIR}/web\"\n    npm install --save-dev ts-protoc-gen\nfi\n\n\ncd ${CURRENT_DIR}\nprotoc -I \"${CURRENT_DIR}/proto\" \"${CURRENT_DIR}/proto/${1}\" --go_out=plugins=grpc:{{.GeneratePath}}\nprotoc -I \"${CURRENT_DIR}/proto\" \\\n    --plugin=\"protoc-gen-ts=${CURRENT_DIR}/web/node_modules/.bin/protoc-gen-ts\" \\\n    --js_out=\"import_style=commonjs,binary:${CURRENT_DIR}/web/src/services\" \\\n    --ts_out=\"${CURRENT_DIR}/web/src/services\" \"${CURRENT_DIR}/proto/${1}\""),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "application/main.go.tmpl",
		FileModTime: time.Unix(1527792315, 0),
		Content:     string("package main\n\nimport (\n\t\"net/http\"\n\n\trice \"github.com/GeertJohan/go.rice\"\n\t\"github.com/gorilla/mux\"\n\t\"github.com/improbable-eng/grpc-web/go/grpcweb\"\n\t\"google.golang.org/grpc\"\n)\n\n//go:generate rice embed-go\n\nfunc main() {\n\tconfig := app.Bootstrap()\n\tmodels.Init(config.DB)\n\n\t{{if .Bootstrap.Assets}}app.Assets = rice.MustFindBox(\"assets\"){{end}}\n\n\tg := grpc.NewServer()\n\tws := grpcweb.WrapServer(g)\n\n\trouter := mux.NewRouter()\n\trouter.PathPrefix(\"/api\").Handler(http.StripPrefix(\"/api\", ws))\n\n\tif app.Env == app.EnvironmentDev {\n\t\trouter.PathPrefix(\"/\").Handler(http.FileServer(http.Dir(\"web/dist\")))\n\t} else {\n\t\tassetsHandler := http.FileServer(rice.MustFindBox(\"web/dist\").HTTPBox())\n\t\trouter.PathPrefix(\"/\").Handler(assetsHandler)\n\t}\n\n\thttp.ListenAndServe(\":\"+config.HTTPPort, router)\n}\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "bootstrap.go.tmpl",
		FileModTime: time.Unix(1527792315, 0),
		Content:     string("package app\n\nimport (\n\t\"database/sql\"\n\t\"log\"\n\t\"os\"\n\n\t\"github.com/joho/godotenv\"\n\t// Load database driver\n\t_ \"github.com/lib/pq\"\n)\n\nconst (\n\t//EnvironmentProd represents production environment\n\tEnvironmentProd = \"PROD\"\n\n\t//EnvironmentDev represents development environment\n\tEnvironmentDev  = \"DEV\"\n)\n\nvar (\n\t// bootstrapped is a flag to prevent multiple bootstrapping\n\tbootstrapped = false\n\n\t// Env indicates in which environment (prod / dev) the application is running\n\tEnv string\n\t{{range .Bootstrap.Settings}}{{if .Public}}\n\t// {{.Name}} {{.Description}}\n\t{{.Name}} string\n\t{{end}}{{end}}\n\t{{if .Bootstrap.Assets}}Assets *rice.Box{{end}}\n)\n\n// Config represents application configuration loaded during bootstrap\ntype Config struct {\n\t{{if not .Bootstrap.NoDB}}DB  *sql.DB{{end}}\n\tHTTPPort string\n\t{{range .Bootstrap.Settings}}{{if not .Public}}\n\t// {{.Name}} {{.Description}}\n\t{{.Name}} string\n\t{{end}}{{end}}\n}\n\n// Bootstrap loads environment variables and initializes the application\nfunc Bootstrap() *Config {\n\tvar config Config\n\n\tif bootstrapped {\n\t\treturn nil\n\t}\n\n\tgodotenv.Load()\n\n\tEnv = os.Getenv(\"ENV\")\n\tif Env == \"\" {\n\t\tEnv = EnvironmentProd\n\t}\n\n\t{{if not .Bootstrap.NoDB}}\n\tdsn := os.Getenv(\"DSN\")\n\tif dsn == \"\" {\n\t\tlog.Fatal(\"Environment variable DSN must be defined. Example: postgres://user:pass@host/db?sslmode=disable\")\n\t}\n\n\tvar err error\n\tconfig.DB, err = sql.Open(\"postgres\", dsn)\n\tif err == nil {\n\t\tlog.Println(\"Connected to database successfully.\")\n\t} else if (Env == EnvironmentDev) {\n\t\tlog.Println(\"Database connection failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database connection failed: \", err)\n\t}\n\n\terr = config.DB.Ping()\n\tif err == nil {\n\t\tlog.Println(\"Pinged database successfully.\")\n\t} else if (Env == EnvironmentDev) {\n\t\tlog.Println(\"Database ping failed: \", err)\n\t} else {\n\t\tlog.Fatal(\"Database ping failed: \", err)\n\t}\n\t{{end}}\n\n\tconfig.HTTPPort = os.Getenv(\"HTTP_PORT\")\n\tif config.HTTPPort == \"\" {\n\t\tconfig.HTTPPort = \"{{.Bootstrap.HTTPPort}}\"\n\t}\n\n\t{{range .Bootstrap.Settings}}{{if not .Public}}\n\tconfig.{{.Name}} = os.Getenv(\"{{upper (snake .Name)}}\")\n\tif config.{{.Name}} == \"\" {\n\t\tlog.Fatal(\"Environment variable {{upper (snake .Name)}} ({{.Description}}) must be defined.\")\n\t}\n\t{{end}}{{end}}\n\n\t{{range .Bootstrap.Settings}}{{if .Public}}\n\t{{.Name}} = os.Getenv(\"{{upper (snake .Name)}}\")\n\tif {{.Name}} == \"\" {\n\t\tlog.Fatal(\"Environment variable {{upper (snake .Name)}} ({{.Description}}) must be defined.\")\n\t}\n\t{{end}}{{end}}\n\n\tos.Clearenv() //prevent non-authorized access\n\n\treturn &config\n}"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "bootstrap_env.tmpl",
		FileModTime: time.Unix(1527792315, 0),
		Content:     string("# The following must be defined as well: ENV{{if not .Bootstrap.NoDB}}, DSN{{end}}, HTTP_PORT\n{{range .Bootstrap.Settings}}{{upper (snake .Name)}} = \"{{.Description}}\"\n{{end}}"),
	}
	file8 := &embedded.EmbeddedFile{
		Filename:    "crud/crud.go.tmpl",
		FileModTime: time.Unix(1527792311, 0),
		Content:     string("package models\n\nimport (\n\t\"context\"\n\t\"database/sql\"\n\t{{range .Imports}}{{.}}{{end}}\n)\n{{.Structure}}\n{{.Get}}\n{{.List}}\n{{.DeleteSingle}}\n{{.DeleteMany}}\n{{.Save}}\n{{.Insert}}\n{{.Update}}\n{{.Merge}}\n{{range .LoadRelated}}{{.}}{{end}}"),
	}
	file9 := &embedded.EmbeddedFile{
		Filename:    "crud/hooks.go.tmpl",
		FileModTime: time.Unix(1527792311, 0),
		Content:     string("package models\nimport (\n\t\"database/sql\"\n)\n\n{{if .Hooks.PreRead}}\nfunc crudPreGet(id {{pkeyPropertyType .Entity.PrimaryKey}}) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PostRead}}\nfunc crudPostGet(entity *{{.Entity.Name}}) error {\n\treturn nil\n}\n{{end}}\n\n{{if .Hooks.PreList}}\nfunc crudPreList(filters []models.ListFilter) ([]models.ListFilter, error) {\n\treturn filters, nil\n}\n{{end}}\n{{if .Hooks.PostList}}\nfunc crudPostList(list []*{{.Entity.Name}}) ([]*{{.Entity.Name}}, error) {\n\treturn list, nil\n}\n{{end}}\n\n{{if .Hooks.PreDelete}}\nfunc crudPreDelete(id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PostDelete}}\nfunc crudPostDelete(id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n\n\n{{if .Hooks.PreSave }}\nfunc crudPreSave(op string, entity *{{.Entity.Name}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n{{if .Hooks.PreSave }}\nfunc crudPostSave(op string, entity *{{.Entity.Name}}, tx *sql.Tx) error {\n\treturn nil\n}\n{{end}}\n\n"),
	}
	filea := &embedded.EmbeddedFile{
		Filename:    "crud/models.go.tmpl",
		FileModTime: time.Unix(1527798334, 0),
		Content:     string("package models\n\nvar (\n\tdb *sql.DB\n\t{{range .Entities}}\n\t{{.Name}}Repo {{.Name}}Repository\n\t{{- end}}\n)\n\nconst (\n\tOperationMerge  byte = 'M'\n\tOperationInsert byte = 'I'\n\tOperationUpdate byte = 'U'\n)\n\n// Init is responsible to initialize all repositories\nfunc Init(database *sql.DB) {\n\tdb = database\n\t{{range .Entities}}\n\t{{.Name}}Repo = {{.Name}}Repositorium{db: database}\n\t{{- end}}\n}\n\n// StartTransaction initiates a database transaction\nfunc StartTransaction() (*sql.Tx, error) {\n\treturn db.Begin()\n}\n\n//ListFilter represents a filter to apply during listing (crud)\ntype ListFilter struct {\n\tField     string\n\tOperation string\n\tValue     interface{}\n}\n\n{{range .Entities -}}\n// {{.Name}}Repository encapsulates operations that may be performed on the entity {{.Name}}\ntype {{.Name}}Repository interface {\n{{if (DerefCrudOpts .Crud).Create -}}\n\t// Insert performs an SQL insert for {{.Name}} record and update instance with inserted id.\n\tInsert(ctx context.Context, entity {{.Name}}, tx *sql.Tx, autocommit bool) ({{.Name}}, error)\n{{end -}}\n{{if (DerefCrudOpts .Crud).Read -}}\n\t// Get returns a single {{.Name}} from database by primary key\n\tGet(ctx context.Context, id {{pkeyPropertyType .PrimaryKey}}) ({{.Name}}, error)\n{{end -}}\n{{if (DerefCrudOpts .Crud).ReadList -}}\n\t// List returns a slice containing {{.Name}} records\n\tList(ctx context.Context, filters []ListFilter, offset, limit int) ([]{{.Name}}, error)\n{{end -}}\n\n{{if (DerefCrudOpts .Crud).Update -}}\n\t// Update Will execute an SQLUpdate Statement for {{.Name}} in the database. Prefer using Save instead of Update directly.\n\tUpdate(ctx context.Context, entity {{.Name}}, tx *sql.Tx, autocommit bool) ({{.Name}}, error)\n{{end -}}\n\n{{if (DerefCrudOpts .Crud).Delete -}}\n\t// DeleteMany deletes many {{.Name}} records from database using filter\n\tDeleteMany(ctx context.Context, filters []ListFilter, tx *sql.Tx, autocommit bool) error\n\t// Delete deletes a {{.Name}} record from database and sets id to nil\n\tDelete(ctx context.Context, entity {{.Name}}, tx *sql.Tx, autocommit bool) ({{.Name}}, error)\n{{end -}}\n{{if (DerefCrudOpts .Crud).Merge -}}\n\t// Merge performs an SQL merge for {{.Name}} record.\n\tMerge(ctx context.Context, entity {{.Name}}, tx *sql.Tx, autocommit bool) ({{.Name}}, error)\n{{end -}}\n{{if (and (and (DerefCrudOpts .Crud).Create (DerefCrudOpts .Crud).Update) (DerefCrudOpts .Crud).Merge) -}}\n\t// Save either inserts or updates a {{.Name}} record based on whether or not id is nil\n\tSave(ctx context.Context, entity {{.Name}}, tx *sql.Tx, autocommit bool) ({{.Name}}, error)\n{{end -}}\n{{- $ThisEntity := .Name -}}\n{{range .Relationships -}}\n{{if eq .Type \"many-many\" -}}\n\t// Load{{.Funcname}} is a helper function to load related {{.Name}} entities\n\tLoad{{RelFuncName .}}(ctx context.Context, entities ...{{$ThisEntity}}) error {\n{{end -}}\n{{if eq .Type \"one-many\" -}}\n\t// Load{{.Funcname}} is a helper function to load related {{.Name}} entities\n\tLoad{{RelFuncName .}}(ctx context.Context, entities ...{{$ThisEntity}}) error {\n{{end -}}\n{{if eq .Type \"many-one\" -}}\n\t// Load{{.Funcname}} is a helper function to load related {{.Name}} entities\n\tLoad{{RelFuncName .}}(ctx context.Context, entities ...{{$ThisEntity}}) error {\n{{end -}}\n{{end -}}\n}\n\n// {{.Name}}Repositorium implements {{.Name}}Repository\ntype {{.Name}}Repositorium struct {\n\tdb *sql.DB\n}\n{{end}}"),
	}
	filec := &embedded.EmbeddedFile{
		Filename:    "crud/partials/delete_many.go.tmpl",
		FileModTime: time.Unix(1527795574, 0),
		Content:     string("\n// DeleteMany deletes many {{.EntityName}} records from database using filter\nfunc (repo {{.EntityName}}Repositorium) DeleteMany(ctx context.Context, filters []ListFilter, tx *sql.Tx, autocommit bool) error {\n\tvar (\n\t\terr      error\n\t\tstmt     *sql.Stmt\n\t\tsegments []string\n\t\tvalues   []interface{}\n\t)\n\n\tif tx == nil {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn err\n\t\t}\n\t\t\n\t\ttx, err = repo.db.Begin()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\tquery := \"DELETE FROM {{.Table}}\"\n\t{{if .HasPreHook}}\n    if filters, err = repo.preDeleteMany(ctx, tx, filters); err != nil {\n\t\treturn err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tfor i, filter := range filters {\n\t\tsegments = append(segments, filter.Field+\" \"+filter.Operation+\" $\"+strconv.Itoa(i+1))\n\t\tvalues = append(values, filter.Value)\n\t}\n\n\tif len(segments) != 0 {\n\t\tquery += \" WHERE \" + strings.Join(segments, \" AND \")\n\t}\n\n\tstmt, err = repo.db.Prepare(query)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\t_, err = stmt.Exec(values...)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tif err = repo.postDeleteMany(ctx, tx, filters); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tif autocommit {\n\t\terr = tx.Commit()\n\t}\n\n\treturn err\n}"),
	}
	filed := &embedded.EmbeddedFile{
		Filename:    "crud/partials/delete_single.go.tmpl",
		FileModTime: time.Unix(1527795594, 0),
		Content:     string("\n// Delete deletes a {{.EntityName}} record from database and sets id to nil\nfunc (repo {{.EntityName}}Repositorium) Delete(ctx context.Context, entity {{.EntityName}}, tx *sql.Tx, autocommit bool) ({{.EntityName}}, error) {\n\tvar (\n\t\terr  error\n\t\tstmt *sql.Stmt\n\t)\n\tid := entity.ID\n\n\tif tx == nil {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn entity, err\n\t\t}\n\n\t\ttx, err = repo.db.Begin()\n\t\tif err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n\tstmt, err = tx.Prepare(\"DELETE FROM {{.Table}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\t{{if .HasPreHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n\tif err = repo.preDelete(ctx, tx, id); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t\n\t_, err = stmt.Exec(id)\n\tif err == nil {\n\t\tentity.ID = {{pkeyPropertyEmptyVal .PrimaryKey}}\n\t} else {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t\n\tif err = repo.postDelete(ctx, tx, id); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\tif autocommit {\n\t\terr = tx.Commit()\n\t}\n\t\n\treturn entity, nil\n}"),
	}
	filee := &embedded.EmbeddedFile{
		Filename:    "crud/partials/get.go.tmpl",
		FileModTime: time.Unix(1527795142, 0),
		Content:     string("\n// Get returns a single {{.EntityName}} from database by primary key\nfunc (repo {{.EntityName}}Repositorium) Get(ctx context.Context, id {{pkeyPropertyType .PrimaryKey}}) ({{.EntityName}}, error) {\n\tvar (\n\t\trows   *sql.Rows\n\t\terr    error\n\t\tentity {{.EntityName}}\n\t)\n\t{{if .HasPreHook}}\n    if err = repo.preGet(ctx, id); err != nil {\n\t\treturn entity, err\n\t}\n    {{end}}\n\t\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n\trows, err = repo.db.Query(\"SELECT {{.SQLFields}} FROM {{.Table}} WHERE id = $1 ORDER BY id ASC\", id)\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\n\tdefer rows.Close()\n\tif rows.Next() {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn entity, err\n\t\t}\n\n\t\t{{range .Before}}{{.}}\n\t\t{{end}}\n\n\t\terr = rows.Scan({{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t\t\n\t\t{{range .After}}{{.}}\n\t\t{{end}}\n\t}\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n\tif err = repo.postGet(ctx, &entity); err != nil {\n\t\treturn entity, err\n\t}\n\t{{end}}\n\n\treturn entity, nil\n}"),
	}
	filef := &embedded.EmbeddedFile{
		Filename:    "crud/partials/insert.go.tmpl",
		FileModTime: time.Unix(1527795466, 0),
		Content:     string("\n// Insert performs an SQL insert for {{.EntityName}} record and update instance with inserted id.\nfunc (repo {{.EntityName}}Repositorium) Insert(ctx context.Context, entity {{.EntityName}}, tx *sql.Tx, autocommit bool) ({{.EntityName}}, error) {\n\tvar (\n\t\t{{- if pkeyIsAuto .PrimaryKey -}}\n\t\tid  {{pkeyPropertyType .PrimaryKey}}\n\t\t{{- end}}\n\t\terr  error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t\t\n\t\ttx, err = repo.db.Begin()\n\t\tif err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t}\n\t{{range .Before}}{{.}}\n\t{{end}}\n\n\t{{if eq .PrimaryKey \"serial\" -}}\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}}) RETURNIentity, NG id\")\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\t{{else}}\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.SQLFields}}) VALUES ({{.SQLPlaceholders}})\")\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\t{{- end}}\n\n\t{{range .After}}{{.}}\n\t{{end}}\n\n\t{{if .HasPreHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\t\n\tif err = repo.preSave(ctx, tx, models.OperationInsert, &entity); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{if eq .PrimaryKey \"serial\" -}}\n\terr = stmt.QueryRow({{.StructFields}}).Scan(&id)\n\tif err == nil {\n\t\tentity.ID = &id\n\t} else {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{else}}\n\t{{if eq .PrimaryKey \"uuid\" -}}\n\tidUUID, err := uuid.NewV4()\n\t\n\tif err == nil {\n\t\tid = idUUID.String()\n\t} else {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\tentity.ID = id\n\t{{- end}}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\tif err := repo.postSave(ctx, \"INSERT\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\tif autocommit {\n\t\terr = tx.Commit()\n\t}\n\n\treturn entity, nil\n}"),
	}
	fileg := &embedded.EmbeddedFile{
		Filename:    "crud/partials/list.go.tmpl",
		FileModTime: time.Unix(1527795522, 0),
		Content:     string("\n// List returns a slice containing {{.EntityName}} records\nfunc (repo {{.EntityName}}Repositorium) List(ctx context.Context, filters []ListFilter, offset, limit int) ([]{{.EntityName}}, error) {\n\tvar (\n\t\tlist     []{{.EntityName}}\n\t\tsegments []string\n\t\tvalues   []interface{}\n\t\terr      error\n\t\trows     *sql.Rows\n\t)\n\n\tquery := \"SELECT {{.SQLFields}} FROM {{.Table}}\"\n\t{{if .HasPreHook}}\n    if filters, err = repo.preList(ctx, filters); err != nil {\n\t\treturn nil, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn nil, err\n\t}\n\n\tfor i, filter := range filters {\n\t\tsegments = append(segments, filter.Field+\" \"+filter.Operation+\" $\"+strconv.Itoa(i+1))\n\t\tvalues = append(values, filter.Value)\n\t}\n\n\tif len(segments) != 0 {\n\t\tquery += \" WHERE \" + strings.Join(segments, \" AND \")\n\t}\n\n\tif limit > -1 {\n\t\tquery += \" LIMIT \"+strconv.Itoa(limit)\n\t}\n\n\tif offset > -1 {\n\t\tquery += \" OFFSET \"+strconv.Itoa(limit)\n\t}\n\n\tquery += \" ORDER BY id ASC\"\n\n\trows, err = repo.db.Query(query, values...)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\tdefer rows.Close()\n\tfor rows.Next() {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\n\t\tvar entity {{.EntityName}}\n\t\t{{range .Before}}{{.}}\n\t\t{{end}}\n\n\t\terr = rows.Scan({{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\t\n\t\t{{range .After}}{{.}}\n\t\t{{end}}\n\n\t\tlist = append(list, entity)\n\t}\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn nil, err\n\t}\n\n\tif list, err = repo.postList(ctx, list); err != nil {\n\t\treturn nil, err\n\t}\n\t{{end}}\n\treturn list, nil\n}"),
	}
	fileh := &embedded.EmbeddedFile{
		Filename:    "crud/partials/loadrelated_manymany.go.tmpl",
		FileModTime: time.Unix(1527795155, 0),
		Content:     string("// Load{{.Funcname}} is a helper function to load related {{.PropertyName}} entities\nfunc (repo {{.ThisEntity}}Repositorium) Load{{.Funcname}}(ctx context.Context, entities ...{{.ThisEntity}}) error {\n\tvar (\n\t\terr error\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices = make(map[{{.ThisType}}][]*{{.ThisEntity}})\n\t)\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tc := 1\n\tfor _, entity := range entities {\n\t\tplaceholder += \"$\" + strconv.Itoa(c) + \",\"\n\t\tindices[entity.ID] = append(indices[entity.ID], &entity)\n\t\tvalues = append(values, entity.ID)\n\t\tc++\n\t}\n\tplaceholder = strings.TrimRight(placeholder, \",\")\n\n\t{{if .Full}}\n\trows, err := repo.db.Query(`\n\t\tSELECT j.{{.ThisID}}, {{.SQLFields}} FROM {{.ThatTable}} t \n\t\tINNER JOIN {{.JoinTable}} j ON t.id = j.{{.ThatID}}\n\t\tWHERE j.{{.ThisID}} IN (`+placeholder+`)\n\t`, values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{else}}\n\trows, err := repo.db.Query(\"SELECT {{.ThisID}}, {{.ThatID}} FROM {{.JoinTable}} WHERE {{.ThisID}} IN (\"+placeholder+\")\", values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{end}}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tfor rows.Next() {\n\t\tvar (\n\t\t\tthisID {{.ThisType}}\n\t\t\t{{if .Full -}}\n\t\t\tentity {{.ThisEntity}}\n\t\t\tthatEntity {{.ThatEntity}}\n\t\t\t{{- else -}}\n\t\t\tthatID {{.ThatType}}\n\t\t\t{{- end -}}\n\t\t)\n\t\t{{if .Full -}}\n\t\t{{range .Before}}{{.}}\n\t\t{{end}}\n\t\terr = rows.Scan(&thisID, {{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t{{range .After}}{{.}}\n\t\t{{end}}\n\t\t{{- else -}}\n\t\terr = rows.Scan(&thisID, &thatID)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t{{- end}}\n\n\t\tfor _, ent := range indices[thisID] {\n\t\t\t{{if .Full -}}\n\t\t\tent.{{.PropertyName}} = append(ent.{{.PropertyName}}, &thatEntity)\n\t\t\t{{- else -}}\n\t\t\tent.{{.PropertyName}} = append(ent.{{.PropertyName}}, thatID)\n\t\t\t{{- end}}\n\t\t}\n\t\t\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\treturn nil\n}"),
	}
	filei := &embedded.EmbeddedFile{
		Filename:    "crud/partials/loadrelated_manyone.go.tmpl",
		FileModTime: time.Unix(1527795164, 0),
		Content:     string("// Load{{.Funcname}} is a helper function to load related {{.PropertyName}} entities\nfunc (repo {{.ThisEntity}}Repositorium) Load{{.Funcname}}(ctx context.Context, entities ...{{.ThisEntity}}) error {\n\tvar (\n\t\terr error\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices = make(map[{{.ThatType}}][]*{{.ThisEntity}})\n\t)\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tc := 1\n\tfor _, entity := range entities {\n\t\tplaceholder += \"$\" + strconv.Itoa(c) + \",\"\n\t\tindices[entity.{{.ThisID}}] = append(indices[entity.{{.ThisID}}], &entity)\n\t\tvalues = append(values, entity.{{.ThisID}})\n\t\tc++\n\t}\n\tplaceholder = strings.TrimRight(placeholder, \",\")\n\trows, err := repo.db.Query(`\n\t\tSELECT id, {{.SQLFields}} FROM {{.ThatTable}} WHERE id IN (`+placeholder+`)\n\t`, values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tfor rows.Next() {\n\t\tvar (\n\t\t\tthatID {{.ThatType}}\n\t\t\tthatEntity {{.ThatEntity}}\n\t\t)\n\t\t{{range .Before}}{{.}}\n\t\t{{end}}\n\t\terr = rows.Scan(&thatID, {{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t{{range .After}}{{.}}\n\t\t{{end}}\n\n\t\tfor _, ent := range indices[thatID] {\n\t\t\tent.{{.PropertyName}} = &thatEntity\n\t\t}\n\t\t\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\treturn nil\n}"),
	}
	filej := &embedded.EmbeddedFile{
		Filename:    "crud/partials/loadrelated_onemany.go.tmpl",
		FileModTime: time.Unix(1527796066, 0),
		Content:     string("// Load{{.Funcname}} is a helper function to load related {{.PropertyName}} entities\nfunc (repo {{.ThisEntity}}Repositorium) Load{{.Funcname}}(ctx context.Context, entities ...{{.ThisEntity}}) error {\n\tvar (\n\t\terr error\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices = make(map[{{.ThisType}}][]*{{.ThisEntity}})\n\t)\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tc := 1\n\tfor _, entity := range entities {\n\t\tplaceholder += \"$\" + strconv.Itoa(c) + \",\"\n\t\tindices[entity.ID] = append(indices[entity.ID], &entity)\n\t\tvalues = append(values, entity.ID)\n\t\tc++\n\t}\n\tplaceholder = strings.TrimRight(placeholder, \",\")\n\n\t{{if .Full}}\n\trows, err := repo.db.Query(`\n\t\tSELECT {{.ThisID}}, {{.SQLFields}} FROM {{.ThatTable}} WHERE {{.ThisID}} IN (`+placeholder+`)\n\t`, values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{else}}\n\trows, err := repo.db.Query(\"SELECT {{.ThisID}}, {{.ThatID}} FROM {{.ThatTable}} WHERE {{.ThisID}} IN (\"+placeholder+\")\", values...)\n\tif err != nil {\n\t\treturn err\n\t}\n\t{{end}}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\n\tfor rows.Next() {\n\t\tvar (\n\t\t\tthisID {{.ThisType}}\n\t\t\t{{if .Full -}}\n\t\t\tentity {{.ThisType}}\n\t\t\tthatEntity {{.ThatEntity}}\n\t\t\t{{- else -}}\n\t\t\tthatID {{.ThatType}}\n\t\t\t{{- end -}}\n\t\t)\n\t\t{{if .Full -}}\n\t\t{{range .Before}}{{.}}\n\t\t{{end}}\n\t\terr = rows.Scan(&thisID, {{.StructFields}})\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t{{range .After}}{{.}}\n\t\t{{end}}\n\t\t{{- else -}}\n\t\terr = rows.Scan(&thisID, &thatID)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\t{{- end}}\n\n\t\tfor _, ent := range indices[thisID] {\n\t\t\t{{if .Full -}}\n\t\t\tent.{{.PropertyName}} = append(ent.{{.PropertyName}}, &thatEntity)\n\t\t\t{{- else -}}\n\t\t\tent.{{.PropertyName}} = append(ent.{{.PropertyName}}, thatID)\n\t\t\t{{- end}}\n\t\t}\n\t\t\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn err\n\t\t}\n\t}\n\n\treturn nil\n}"),
	}
	filek := &embedded.EmbeddedFile{
		Filename:    "crud/partials/merge.go.tmpl",
		FileModTime: time.Unix(1527795661, 0),
		Content:     string("\n// Merge performs an SQL merge for {{.EntityName}} record.\nfunc (repo {{.EntityName}}Repositorium) Merge(ctx context.Context, entity {{.EntityName}}, tx *sql.Tx, autocommit bool) ({{.EntityName}}, error) {\n\tvar (\n\t\terr error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t\t\n\t\ttx, err = repo.db.Begin()\n\t\tif err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t}\n\n\tif entity.ID == {{pkeyPropertyEmptyVal .PrimaryKey}} {\n\t\treturn {{.EntityName}}Repo.Insert(ctx, entity, tx, autocommit)\n\t}\n\n\t{{range .Before}}{{.}}\n\t{{end}}\n\n\tstmt, err = tx.Prepare(`INSERT INTO {{.Table}} ({{.SQLFieldsInsert}}) VALUES ({{.SQLPlaceholders}}) \n\tON CONFLICT (id) DO UPDATE SET {{.SQLFieldsUpdate}}`)\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\t{{if .HasPreHook}}\n    if err = repo.preSave(ctx, tx, models.OperationMerge, &entity); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\t{{range .After}}{{.}}\n\t{{end}}\n\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\tif err = repo.postSave(ctx, \"MERGE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n\tif autocommit {\n\t\terr = tx.Commit()\n\t}\n\n\treturn entity, err\n}"),
	}
	filel := &embedded.EmbeddedFile{
		Filename:    "crud/partials/save.go.tmpl",
		FileModTime: time.Unix(1527795685, 0),
		Content:     string("\n// Save either inserts or updates a {{.EntityName}} record based on whether or not id is nil\nfunc (repo {{.EntityName}}Repositorium) Save(ctx context.Context, entity {{.EntityName}}, tx *sql.Tx, autocommit bool) ({{.EntityName}}, error) {\n\t{{if pkeyIsAuto .PrimaryKey -}}\n\tif entity.ID == {{pkeyPropertyEmptyVal .PrimaryKey}} {\n\t\treturn {{.EntityName}}Repo.Insert(ctx, entity, tx, autocommit)\n\t}\n\treturn {{.EntityName}}Repo.Update(ctx, entity, tx, autocommit)\n\t{{- else -}}\n\tif entity.ID == {{pkeyPropertyEmptyVal .PrimaryKey}} {\n\t\treturn entity, errors.New(\"primary key cannot be nil\")\n\t}\n\treturn {{.EntityName}}Repo.Merge(ctx, entity, tx, autocommit)\n\t{{end -}}\n}"),
	}
	filem := &embedded.EmbeddedFile{
		Filename:    "crud/partials/saverelated.go.tmpl",
		FileModTime: time.Unix(1527792701, 0),
		Content:     string("// SaveRelated{{.PropertyName}} is a helper function to save related {{.PropertyName}} entities\nfunc SaveRelated{{.PropertyName}}(ctx context.Context, idthis  {{pkeyPropertyType .PrimaryKey}}, relatedids {{.PropertyType}}) error {\n\tvar (\n\t\tplaceholder string\n\t\tvalues  []interface{}\n\t\tindices map[{{pkeyPropertyType .PrimaryKey}}]{{trimPrefix .PropertyType \"[]\"}}\n\t\tidthis  {{pkeyPropertyType .PrimaryKey}}\n\t\tidthat  {{trimPrefix .PropertyType \"[]\"}}\n\t\titems   []{{.EntityName}}\n\t\tstmt    *sql.Stmt\n\t)\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn err\n\t}\n\t\n\tstmt, err = tx.Prepare(\"DELETE FROM {{.Table}} WHERE {{.ThisID}} = $1\")\n\n\tif err != nil {\n\t\treturn err\n\t}\n\n\t_, err = stmt.Exec(entity.ID)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn err\n\t}\n\n\tstmt, err = tx.Prepare(\"INSERT INTO {{.Table}} ({{.ThisID}}, {{.ThatID}}) VALUES ($1, $2)\")\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tfor _, relatedID := range *entity.{{.PropertyName}} {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\n\t\t_, err = stmt.Exec(entity.ID, relatedID)\n\t\tif err != nil {\n\t\t\ttx.Rollback()\n\t\t\treturn err\n\t\t}\n\t}\n}"),
	}
	filen := &embedded.EmbeddedFile{
		Filename:    "crud/partials/update.go.tmpl",
		FileModTime: time.Unix(1527795142, 0),
		Content:     string("\n// Update Will execute an SQLUpdate Statement for {{.EntityName}} in the database. Prefer using Save instead of Update directly.\nfunc (repo {{.EntityName}}Repositorium) Update(ctx context.Context, entity {{.EntityName}}, tx *sql.Tx, autocommit bool) ({{.EntityName}}, error) {\n\tvar (\n\t\terr error\n\t\tstmt *sql.Stmt\n\t)\n\n\tif tx == nil {\n\t\tif err = util.CheckContext(ctx); err != nil {\n\t\t\treturn entity, err\n\t\t}\n\n\t\ttx, err = repo.db.Begin()\n\t\tif err != nil {\n\t\t\treturn entity, err\n\t\t}\n\t}\n\t\n\t{{range .Before}}{{.}}\n\t{{end}}\n\n\tstmt, err = tx.Prepare(\"UPDATE {{.Table}} SET {{.SQLFields}} WHERE id = $1\")\n\tif err != nil {\n\t\treturn entity, err\n\t}\n\n\t{{range .After}}{{.}}\n\t{{end}}\n\n\t{{if .HasPreHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\treturn entity, err\n\t}\n\n    if err = repo.preSave(ctx, tx, models.OperationUpdate, &entity); err != nil {\n\t\ttx.Rollback()\n        return entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t_, err = stmt.Exec({{.StructFields}})\n\tif err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{if .HasPostHook}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\tif err = repo.postSave(ctx, \"UPDATE\", entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\t{{end}}\n\tif err = util.CheckContext(ctx); err != nil {\n\t\ttx.Rollback()\n\t\treturn entity, err\n\t}\n\n\tif autocommit {\n\t\terr = tx.Commit()\n\t}\n\n\treturn entity, err\n}"),
	}
	fileo := &embedded.EmbeddedFile{
		Filename:    "crud/protobuf.proto.tmpl",
		FileModTime: time.Unix(1527839828, 0),
		Content:     string("syntax = \"proto3\";\n\npackage models;\noption go_package = \"{{.ProjectImportPath}}/models\";\n\n{{range .Imports -}}\n{{.}}\n{{- end}}\n\n{{ range .Entities -}}\n// {{.Name}} {{.Description}}\nmessage {{.Name}} { {{ range .Fields }}\n\t{{.Type}} {{.Name}} = {{.Index}};\n{{- end}}\n}\n\n{{end}}"),
	}
	filep := &embedded.EmbeddedFile{
		Filename:    "http.go.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("package app\n\nimport (\n\t\"log\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"syscall\"\n\n\t\"github.com/gorilla/mux\"\n)\n\n// ServeHTTP starts an http server\nfunc ServeHTTP(listen string, route func(router *mux.Router) error) {\n\tvar err error\n\tsigs := make(chan os.Signal, 1)\n\tsignal.Notify(sigs, syscall.SIGTERM)\n\n\trouter := mux.NewRouter()\n\terr = route(router)\n\n\tif err != nil {\n\t\tlog.Fatal(\"Failed to register routes: \", err)\n\t}\n\n\tgo func() {\n\t\terr = http.ListenAndServe(listen, router)\n\t\tif err != nil {\n\t\t\tlog.Fatal(\"Failed to start http server: \", err)\n\t\t}\n\t}()\n\n\tlog.Println(\"Listening on \" + listen)\n\t<-sigs\n\tlog.Println(\"Server stopped\")\n}\n"),
	}
	fileq := &embedded.EmbeddedFile{
		Filename:    "rest.go.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("package {{.Package}}\nimport (\n\t\"database/sql\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"net/http\"\n\t\"strconv\"\n\n\t\"github.com/gorilla/mux\"\n)\n\ntype responseSingle struct {\n\tStatus   bool      `json:\"status\"`\n\tMessages []message `json:\"messages\"`\n\tEntity   *{{.Entity.Name}} `json:\"entity\"`\n}\n\ntype responseList struct {\n\tStatus   bool                `json:\"status\"`\n\tMessages []message           `json:\"messages\"`\n\tEntities []*{{.Entity.Name}} `json:\"entities\"`\n}\n\ntype message struct {\n\tType    rune   `json:\"type\"`\n\tMessage string `json:\"message\"`\n}\n\n//RegisterRoutes registers routes with a mux Router\nfunc RegisterRoutes(router *mux.Router) {\n\t{{if .Entity.Rest.Read}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestGet).Methods(\"GET\"){{end}}\n\t{{if .Entity.Rest.ReadList}}router.HandleFunc(\"/{{.Endpoint}}\", RestList).Methods(\"GET\"){{end}}\n\t{{if .Entity.Rest.Create}}router.HandleFunc(\"/{{.Endpoint}}\", RestCreate).Methods(\"POST\"){{end}}\n\t{{if .Entity.Rest.Update}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestUpdate).Methods(\"PUT\"){{end}}\n\t{{if .Entity.Rest.Delete}}router.HandleFunc(\"/{{.Endpoint}}/{id}\", RestDelete).Methods(\"DELETE\"){{end}}\n}\n\n{{if .Entity.Rest.Read}}\n//RestGet is a REST endpoint for GET /{{.Endpoint}}/{id}\nfunc RestGet(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\terr      error\n\t\tresponse responseSingle\n\t\t{{if or .Entity.Rest.Hooks.PreRead .Entity.Rest.Hooks.PostRead -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreRead}}\n    if stop, err = restPreGet(w, r, id); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostRead}}\n    if stop, err = restPostGet(w, r, response.Entity); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Status = true\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.ReadList}}\n//RestList is a REST endpoint for GET /{{.Endpoint}}\nfunc RestList(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\tresponse responseList\n\t\tfilters  []models.ListFilter\n\t\t{{if or .Entity.Rest.Hooks.PreList .Entity.Rest.Hooks.PostList}}stop     bool{{end}}\n\t)\n\t{{range .Entity.Fields}}{{if .Filterable}}\n\t{{if eq .Property.Type \"bool\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); val != \"\" {\n\t\tif t, e := strconv.ParseBool(val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:t})\n\t\t}\n\t}\n\t{{end}}\n\t{{if eq .Property.Type \"string\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); val != \"\" {\n\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:val})\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-lk\"); val != \"\" {\n\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"LIKE\", Value:\"%\" + val + \"%\"})\n\t}\n\t{{end}}\n\t{{if eq .Property.Type \"time.Time\"}}\n\tif val := query.Get(\"{{.Serialized}}\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"=\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-gt\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\">\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-ge\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\">=\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-lt\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"<\", Value:t})\n\t\t}\n\t}\n\n\tif val := query.Get(\"{{.Serialized}}-le\"); len(val) == 16 {\n\t\tif t, e := time.Parse(\"2006-01-02-15-04\", val); e == nil {\n\t\t\tfilters = append(filters, models.ListFilter{Field:\"{{.Schema.Field}}\", Operation:\"<=\", Value:t})\n\t\t}\n\t}\n\t{{end}}\n\t{{end}}{{end}}\n\n\t{{if .Entity.Rest.Hooks.PreList}}\n    if filters, stop, err = restPreList(w, r, filters); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Entities, err = List(filters)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostList}}\n    if response.Entities, stop, err = restPostList(w, r, response.Entities); err != nil || stop {\n        return\n    }\n    {{end}}\n\n\tresponse.Status = true\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Create}}\n//RestCreate is a REST endpoint for POST /{{.Endpoint}}\nfunc RestCreate(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\trawbody  []byte\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreCreate .Entity.Rest.Hooks.PostCreate}}stop     bool{{end}}\n\t)\n\n\trawbody, err = ioutil.ReadAll(r.Body)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to read body\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity = New()\n\terr = json.Unmarshal(rawbody, response.Entity)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to decode body\"}]}`)\n\t\treturn\n\t}\n\t{{if pkeyIsAuto .Entity.PrimaryKey -}}\n\tresponse.Entity.ID = nil\n\t{{- end}}\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreCreate}}\n\tif stop, err = restPreCreate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\n\terr = response.Entity.Save(tx, false)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Save failed\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostCreate}}\n\tif stop, err = restPostCreate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\t\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestCreate could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Update}}\n//RestUpdate is a REST endpoint for PUT /{{.Endpoint}}/{id}\nfunc RestUpdate(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\terr      error\n\t\trawbody  []byte\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreUpdate .Entity.Rest.Hooks.PostUpdate -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\trawbody, err = ioutil.ReadAll(r.Body)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to read body\"}]}`)\n\t\treturn\n\t}\n\n\terr = json.Unmarshal(rawbody, response.Entity)\n\tif err != nil {\n\t\tif err != nil {\n\t\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\t\tw.WriteHeader(http.StatusBadRequest)\n\t\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to decode body\"}]}`)\n\t\t\treturn\n\t\t}\n\t}\n\tresponse.Entity.ID = &id\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PreUpdate}}\n    if stop, err = restPreUpdate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n        return\n    } else if stop {\n\t\treturn\n\t}\n    {{end}}\n\n\terr = response.Entity.Save(tx, false)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Save failed\"}]}`)\n\t\treturn\n\t}\n\n\t{{if .Entity.Rest.Hooks.PostUpdate}}\n    if stop, err = restPostUpdate(w, r, response.Entity, tx); err != nil {\n\t\ttx.Rollback()\n        return\n    } else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\t\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestUpdate could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n\n{{if .Entity.Rest.Delete}}\n//RestDelete is a REST endpoint for DELETE /{{.Endpoint}}/{id}\nfunc RestDelete(w http.ResponseWriter, r *http.Request) {\n\tvar (\n\t\tid       {{pkeyPropertyType .Entity.PrimaryKey}}\n\t\terr      error\n\t\tresponse responseSingle\n\t\ttx       *sql.Tx\n\t\t{{if or .Entity.Rest.Hooks.PreDelete .Entity.Rest.Hooks.PostDelete -}}\n\t\tstop     bool\n\t\t{{- end}}\n\t)\n\n\tvars := mux.Vars(r)\n\t{{if pkeyIsInt .Entity.PrimaryKey -}}\n\tvalid := false\n\tif _, ok := vars[\"id\"]; ok {\n\t\tid, err = strconv.ParseInt(vars[\"id\"], 10, 64)\n\t\tvalid = err == nil && id > 0\n\t}\n\t{{else}}\n\tid, valid := vars[\"id\"]\n\t{{- end}}\n\n\tif !valid {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Invalid ID\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Entity, err = Get(id)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"An error occurred\"}]}`)\n\t\treturn\n\t}\n\n\tif response.Entity == nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusNotFound)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Entity not found\"}]}`)\n\t\treturn\n\t}\n\n\ttx, err = db.Begin()\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Failed to process\"}]}`)\n\t\treturn\n\t}\n\t{{if .Entity.Rest.Hooks.PreDelete}}\n\tif stop, err = restPreDelete(w, r, id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n    {{end}}\n\terr = response.Entity.Delete(tx, false)\n\tif err != nil {\n\t\ttx.Rollback()\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"Delete failed\"}]}`)\n\t\treturn\n\t}\n\t{{if .Entity.Rest.Hooks.PostDelete}}\n\tif stop, err = restPostDelete(w, r, id, tx); err != nil {\n\t\ttx.Rollback()\n\t\treturn\n\t} else if stop {\n\t\treturn\n\t}\n\t{{end}}\n\tif err = tx.Commit(); err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusBadRequest)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"E\", \"message\": \"RestDelete could not commit transaction\"}]}`)\n\t\treturn\n\t}\n\n\tresponse.Status = true\t\n\toutput, err := json.Marshal(response)\n\tif err != nil {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\tw.WriteHeader(http.StatusInternalServerError)\n\t\tfmt.Fprint(w, `{\"status\": false, \"messages\": [{\"type\": \"error\", \"text\": \"JSON encoding failed\"}]}`)\n\t\treturn\n\t}\n\n\tw.Header().Set(\"Content-Type\", \"application/json\")\n\tw.WriteHeader(http.StatusOK)\n\tfmt.Fprint(w, string(output))\n}\n{{end}}\n"),
	}
	filer := &embedded.EmbeddedFile{
		Filename:    "rest_hooks.go.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("package {{.Package}}\n\nimport (\n\t\"database/sql\"\n\t\"net/http\"\n)\n\n{{if .Hooks.PreRead}}\nfunc restPreGet(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostRead}}\nfunc restPostGet(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreList}}\nfunc restPreList(w http.ResponseWriter, r *http.Request, filters []models.ListFilter) ([]models.ListFilter, bool, error) {\n\treturn filters, false, nil\n}\n{{end}}\n{{if .Hooks.PostList}}\nfunc restPostList(w http.ResponseWriter, r *http.Request, list []*{{.Entity.Name}}) ([]*{{.Entity.Name}}, bool, error) {\n\treturn list, false, nil\n}\n{{end}}\n\n{{if .Hooks.PreCreate}}\nfunc restPreCreate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostCreate}}\nfunc restPostCreate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreUpdate}}\nfunc restPreUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostUpdate}}\nfunc restPostUpdate(w http.ResponseWriter, r *http.Request, entity *{{.Entity.Name}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n\n{{if .Hooks.PreDelete}}\nfunc restPreDelete(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}\n{{if .Hooks.PostDelete}}\nfunc restPostDelete(w http.ResponseWriter, r *http.Request, id {{pkeyPropertyType .Entity.PrimaryKey}}, tx *sql.Tx) (bool, error) {\n\treturn false, nil\n}\n{{end}}"),
	}
	files := &embedded.EmbeddedFile{
		Filename:    "schema.sql.tmpl",
		FileModTime: time.Unix(1527792315, 0),
		Content:     string("DROP TABLE IF EXISTS {{.Entity.Table}};\n\nCREATE TABLE {{.Entity.Table}} (\n\t\"id\" {{pkeyFieldType .Entity.PrimaryKey}},\n\t{{- range $i, $e := .Entity.Fields}}{{if ne .Schema.Field \"\"}}\n\t\"{{.Schema.Field}}\" {{$e.Schema.Type}} NOT NULL\n\t{{- if ne .Schema.Default \"\" -}} DEFAULT {{.Schema.Default}}{{end}},\n\t{{- end}}{{- end}}\n\t{{- range .RelatedFields}}\n\t\"{{.Name}}\" {{.Type}} NOT NULL,\n\t{{end}}\n\t{{range .Entity.TableConstraints}}{{.}},{{end}}\n\tPRIMARY KEY (\"id\")\n);\n\n{{range .RelatedTables}}\nDROP TABLE IF EXISTS {{.Table}};\n\nCREATE TABLE {{.Table}} (\n\t\"{{.ThisID}}\" {{.ThisType}} NOT NULL,\n\t\"{{.ThatID}}\" {{.ThatType}} NOT NULL\n);\n{{end}}"),
	}
	fileu := &embedded.EmbeddedFile{
		Filename:    "util/util.go.tmpl",
		FileModTime: time.Unix(1527792315, 0),
		Content:     string("package util\n\nimport (\n\t\"context\"\n)\n\n// CheckContext returns an error if context is done\nfunc CheckContext(ctx context.Context) error {\n\tselect {\n\tcase <-ctx.Done():\n\t\treturn ctx.Err()\n\tdefault:\n\t\treturn nil\n\t}\n}"),
	}
	filev := &embedded.EmbeddedFile{
		Filename:    "vuetify_actions.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("import types from \"./types\";\n\nexport default {}"),
	}
	filew := &embedded.EmbeddedFile{
		Filename:    "vuetify_edit.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<template>\n    <div class=\"container\">\n\t\t<v-toolbar color=\"transparent\" flat>\n            <v-toolbar-title class=\"grey--text text--darken-4 ml-0\"><h2>{{.Entity.Name}}</h2></v-toolbar-title>\n            <v-spacer></v-spacer>\n            <v-btn ml-0 small color=\"grey\" flat :to=\"{name: '{{.Endpoint}}List'}\">\n                <v-icon dark>arrow_back</v-icon> Back\n            </v-btn>\n        </v-toolbar>\n\t\t<v-alert :type=\"message.type\" :value=\"true\" v-for=\"(message, index) in messages\" :key=\"index\">\n\t\t{{ \"{{ message.text }}\" }}\n\t\t</v-alert>\n  \n        {{range .Entity.Fields -}}\n        {{widget_field \"vuetify\" .Widget.Type .}}\n        {{- end}}\n\n        <v-btn color=\"primary\" @click=\"save()\">Save</v-btn>\n        <v-btn color=\"gray\" :to=\"{name: '{{.Endpoint}}List'}\">Cancel</v-btn>\n\t</div>\n</template>\n  \n<script>\nimport axios from \"axios\"\n\nexport default {\n    props: [\"id\"],\n    created() {\n        if (!this.id) {\n            return\n        }\n\n        axios.get(\"/api/{{.Endpoint}}/\" + this.id).then(response => {\n            this.id = response.data.entity.id\n            {{range .Entity.Fields}}{{if ne .Serialized \"id\"}}\n            this.entity.{{.Serialized}} = response.data.entity.{{.Serialized}}\n            {{if eq .Widget.Type \"date\"}}this.dates.{{.Serialized}}.value = response.data.entity.{{.Serialized}}.substr(0,10){{end}}\n            {{end}}{{end}}\n        })\n    },\n    data() {\n        return {\n            select: {\n                {{range $i, $v := .Entity.Fields}}{{if eq .Widget.Type \"date\"}}\n                {{.Serialized}}: {\n                    items:[\n                        {{range $j, $u := .Widget.Options}}\n                        {text: \"{{.Label}}\", value: \"{{.Value}}\"}{{if eq (plus1 $j) (len $u)}},{{end}}\n                        {{end}}\n                    ]\n                }{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n\t\t\t},\n\t\t\tdates: {\n                {{range $i, $v := .Entity.Fields}}{{if eq .Widget.Type \"date\"}}\n\t\t\t\t{{.Serialized}}: {value: null, menuAppear: false}{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n\t\t\t},\n            messages: [],\n            entity: {\n                {{range $i, $e := .Entity.Fields}}{{if ne .Serialized \"id\"}}\n                {{.Serialized}} : null{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n                {{end}}{{end}}\n            }\n        }\n    },\n    watch: {\n        {{range $i, $e := .Entity.Fields}}\n        \"select.{{.Serialized}}.search\": function(val) {\n            val && this.querySelections(\"{{.Serialized}}\", \"{{$.Endpoint}}\", \"{{$.Prefix}}{{.Relationship.Target.Endpoint}}\", val)\n        }{{if ne (plus1 $i) (len $.Entity.Fields)}},{{end}}\n        {{end}}\n    },\n    methods: {\n        querySelections(fieldname, endpoint, filter, val) {\n            this.select[fieldname].loading = true\n            axios.get(\"/api/\" + endpoint + \"?\" + filter + \"-lk=\" + encodeURIComponent(val)).then(response => {\n                this.select[fieldname].loading = false\n                this.select[fieldname].items = response.data.entities.map(function(e) {\n                    return { text: e[filter], value: e.id }\n                })\n            })\n        },\n        save() {\n            if (this.id) {\n                axios.put(\"/api/{{.Endpoint}}/\" + this.id, this.entity).then(this.saved)\n            } else {\n                axios.post(\"/api/{{.Endpoint}}\", this.entity).then(this.saved)\n            }\n\t\t},\n\t\tsaved(response) {\n\t\t\tthis.id = response.data.entity.id\n\t\t\t{{range .Entity.Fields}}{{if ne .Serialized \"id\"}}\n            this.entity.{{.Serialized}} = response.data.entity.{{.Serialized}}\n            {{end}}{{end}}\n\n\t\t\tthis.messages.push({\n\t\t\t\ttype: \"success\",\n\t\t\t\ttext: \"{{.Entity.Name}} saved successfully\"\n\t\t\t})\n\t\t}\n    }\n}\n</script>"),
	}
	filex := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-checkbox.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-checkbox label=\"{{.Label}}\" v-model=\"entity.{{.Serialized}}\"></v-checkbox>"),
	}
	filey := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-date.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-menu\n\tref=\"menu_{{.Serialized}}\"\n\tlazy\n\t:close-on-content-click=\"false\"\n\tv-model=\"dates.{{.Serialized}}.menuAppear\"\n\ttransition=\"scale-transition\"\n\toffset-y\n\tfull-width\n\t:nudge-right=\"40\"\n\tmin-width=\"290px\"\n\t:return-value.sync=\"dates.{{.Serialized}}.value\"\n\t>\n\t<v-text-field\n\t\tslot=\"activator\"\n\t\tlabel=\"{{.Label}}\"\n\t\tv-model=\"dates.{{.Serialized}}.value\"\n\t\tprepend-icon=\"event\"\n\t\treadonly\n\t\t></v-text-field>\n\t\t<v-date-picker v-model=\"dates.{{.Serialized}}.value\" @change=\"entity.{{.Serialized}} = dates.{{.Serialized}}.value + 'T00:00:00Z'\" no-title scrollable>\n\t\t<v-spacer></v-spacer>\n\t\t<v-btn flat color=\"primary\" @click=\"menu_{{.Serialized}} = false\">Cancel</v-btn>\n\t\t<v-btn flat color=\"primary\" @click=\"$refs.menu_{{.Serialized}}.save(dates.{{.Serialized}}.value)\">OK</v-btn>\n\t\t</v-date-picker>\n</v-menu>"),
	}
	filez := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-number.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" type=\"number\" />"),
	}
	file10 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-password.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-text-field\n\tv-model=\"entity.{{.Serialized}}\"\n\t:append-icon=\"e1 ? 'visibility' : 'visibility_off'\"\n\t:append-icon-cb=\"() => (e1 = !e1)\"\n\t:type=\"e1 ? 'password' : 'text'\"\n\tcounter\n  ></v-text-field>"),
	}
	file11 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-select-rel.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-select\n    autocomplete\n    cache-items\n    required\n    label=\"{{.Label}}\"\n    :loading=\"select.{{.Serialized}}.isloading\"\n    :items=\"select.{{.Serialized}}.items\"\n\t:search-input.sync=\"select.{{.Serialized}}.search\"\n\tv-model=\"entity.{{.Serialized}}\"\n\t{{if .Widget.Multiple}}multiple chips{{end}}\n></v-select>"),
	}
	file12 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-select.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-select\n\tautocomplete\n\tcache-items\n\trequired\n\tlabel=\"{{.Label}}\"\n\t:items=\"select.{{.Serialized}}.items\"\n\tv-model=\"entity.{{.Serialized}}\"\n\t{{if .Widget.Multiple}}multiple chips{{end}}\n></v-select>"),
	}
	file13 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-textarea.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" multiline />"),
	}
	file14 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-textfield.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-text-field v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" />"),
	}
	file15 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-time.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<div>\n\t<v-time-picker v-model=\"entity.{{.Serialized}}\" label=\"{{.Label}}\" :landscape=\"landscape\"></v-time-picker>\n</div>"),
	}
	file16 := &embedded.EmbeddedFile{
		Filename:    "vuetify_editor-field-toggle.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<v-switch\n\tlabel=\"{{.Label}}\"\n\tv-model=\"entity.{{.Serialized}}\"\n></v-switch>"),
	}
	file17 := &embedded.EmbeddedFile{
		Filename:    "vuetify_getters.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("export default {}\n"),
	}
	file18 := &embedded.EmbeddedFile{
		Filename:    "vuetify_index.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("import actions from \"./actions\";\nimport getters from \"./getters\";\nimport mutations from \"./mutations\";\nimport routes from \"./routes\";\n\nconst namespaced = true;\n\nconst state = {\n  entities: routes.routes\n};\n\nexport default {\n  namespaced,\n  state,\n  actions,\n  getters,\n  mutations\n};"),
	}
	file19 := &embedded.EmbeddedFile{
		Filename:    "vuetify_list.vue.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("<template>\n    <v-container>\n        <v-toolbar color=\"transparent\" flat>\n            <v-toolbar-title class=\"grey--text text--darken-4 ml-0\"><h2>{{.Entity.Name}}</h2></v-toolbar-title>\n            <v-spacer></v-spacer>\n            <v-btn mr-0 color=\"primary\" :to=\"{name: '{{.Endpoint}}Edit', params:{id: 0}}\">\n                <v-icon dark>add</v-icon> Add\n            </v-btn>\n        </v-toolbar>\n        \n        <v-alert :type=\"message.type === 'E' ? 'error' : message.type\" :value=\"true\" v-for=\"(message, index) in messages\" :key=\"index\">\n            {{ \"{{ message.text }}\" }}\n        </v-alert>\n\n        <v-alert type=\"info\" value=\"true\"  color=\"primary\" outline icon=\"info\" v-if=\"entities.length === 0\">\n            No {{.Entity.Name}} exist. Would you like to create one now?\n            <v-btn :to=\"{name: '{{.Endpoint}}Edit', params:{id: 0}}\" color=\"primary\">create new</v-btn>\n        </v-alert>\n        <template v-else>\n            <v-text-field mb-4 append-icon=\"search\" label=\"Search\" single-line hide-details v-model=\"search\"></v-text-field>            \n            <v-data-table :headers=\"headers\" :items=\"entities\" class=\"elevation-1\" :search=\"search\">\n                <template slot=\"items\" slot-scope=\"props\">\n\t\t\t\t\t{{ range .Entity.Fields }}\n\t\t\t\t\t<td>{{ printf \"{{ props.item.%s}}\" .Serialized }}</td>\n\t\t\t\t\t{{end}}\n                    <td class=\"justify-center layout px-0\">\n                        <v-btn icon class=\"mx-0\" :to=\"{name: '{{.Endpoint}}Edit', params: {'id': props.item.id}  }\">\n                            <v-icon color=\"teal\">edit</v-icon>\n                        </v-btn>\n                    </td>\n                </template>\n\n                <template slot=\"no-data\">\n                    <v-flex ma-4>\n                        <v-alert slot=\"no-results\" :value=\"true\" color=\"primary\" outline icon=\"info\" v-if=\"search.length > 0\">\n                        Your search for \"{{ \"{{ search }}\" }}\" found no results.\n                        </v-alert>\n                        <v-alert slot=\"no-results\" :value=\"true\" color=\"primary\" outline icon=\"info\" v-else>\n                            No {{.Entity.Name}} found.\n                        </v-alert>\n                    </v-flex>\n                </template>\n            </v-data-table>\n        </template>\n    </v-container>\n</template>\n\n<script>\nimport axios from \"axios\"\nexport default {\n  data() {\n    return {\n      messages: [],\n      search: \"\",\n      headers: [\n\t\t{{range .Entity.Fields }}\n\t\t{text: \"{{.Label}}\", value: \"{{.Serialized}}\"},\n\t\t{{end}}\n        {'text': 'Action', 'value': null}\n      ],\n      entities: []\n    };\n  },\n  created() {\n    axios\n      .get(\"/api/{{.Endpoint}}\")\n      .then(response => {\n        this.entities = response.data.entities;\n      })\n      .catch(error => {\n        this.messages = [...this.messages, ...error.response.data.messages];\n      });\n  }\n};\n</script>"),
	}
	file1a := &embedded.EmbeddedFile{
		Filename:    "vuetify_mutations.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("import types from \"./types\";\n\nexport default {}\n"),
	}
	file1b := &embedded.EmbeddedFile{
		Filename:    "vuetify_routes.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("{{range .Entities}}\n// {{.Name}} {{.Description}}\nimport {{.Name}}Edit from \"../views/{{plural .Name}}Edit.vue\";\nimport {{.Name}}List from \"../views/{{plural .Name}}List.vue\";\n{{end}}\n\nlet routes = [\n  {{range $i, $v := .Entities}}\n  {\n    path: \"/{{lower (plural .Name)}}/:id\",\n    name: \"{{lower (plural .Name)}}Edit\",\n    props: true,\n    icon: \"dashboard\",\n    component: {{.Name}}Edit,\n    entity: \"{{plural .Name}}\"\n  },\n  {\n    path: \"/{{lower (plural .Name)}}list/\",\n    name: \"{{lower (plural .Name)}}List\",\n    icon: \"dashboard\",\n    component: {{.Name}}List,\n    entity: \"{{plural .Name}}\"\n  }{{if ne (plus1 $i) (len $.Entities)}},{{end}}\n  {{end}}\n];\n\nlet entities = [\n  {{range $i, $v := .Entities}}\n  \"{{plural .Name}}\"{{if ne (plus1 $i) (len $.Entities)}},{{end}}\n  {{end}}\n];\n\nfunction registerRoutes(router) {\n  router.addRoutes(routes);\n}\n\nexport default {\n  routes,\n  entities,\n  registerRoutes\n}\n"),
	}
	file1c := &embedded.EmbeddedFile{
		Filename:    "vuetify_types.js.tmpl",
		FileModTime: time.Unix(1527597130, 0),
		Content:     string("export default {}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1527792315, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file5,  // "bootstrap.go.tmpl"
			file6,  // "bootstrap_env.tmpl"
			filep,  // "http.go.tmpl"
			fileq,  // "rest.go.tmpl"
			filer,  // "rest_hooks.go.tmpl"
			files,  // "schema.sql.tmpl"
			filev,  // "vuetify_actions.js.tmpl"
			filew,  // "vuetify_edit.vue.tmpl"
			filex,  // "vuetify_editor-field-checkbox.vue.tmpl"
			filey,  // "vuetify_editor-field-date.vue.tmpl"
			filez,  // "vuetify_editor-field-number.vue.tmpl"
			file10, // "vuetify_editor-field-password.vue.tmpl"
			file11, // "vuetify_editor-field-select-rel.vue.tmpl"
			file12, // "vuetify_editor-field-select.vue.tmpl"
			file13, // "vuetify_editor-field-textarea.vue.tmpl"
			file14, // "vuetify_editor-field-textfield.vue.tmpl"
			file15, // "vuetify_editor-field-time.vue.tmpl"
			file16, // "vuetify_editor-field-toggle.vue.tmpl"
			file17, // "vuetify_getters.js.tmpl"
			file18, // "vuetify_index.js.tmpl"
			file19, // "vuetify_list.vue.tmpl"
			file1a, // "vuetify_mutations.js.tmpl"
			file1b, // "vuetify_routes.js.tmpl"
			file1c, // "vuetify_types.js.tmpl"

		},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "application",
		DirModTime: time.Unix(1527792315, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "application/gen-service.sh.tmpl"
			file4, // "application/main.go.tmpl"

		},
	}
	dir7 := &embedded.EmbeddedDir{
		Filename:   "crud",
		DirModTime: time.Unix(1527798970, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file8, // "crud/crud.go.tmpl"
			file9, // "crud/hooks.go.tmpl"
			filea, // "crud/models.go.tmpl"
			fileo, // "crud/protobuf.proto.tmpl"

		},
	}
	dirb := &embedded.EmbeddedDir{
		Filename:   "crud/partials",
		DirModTime: time.Unix(1527792315, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			filec, // "crud/partials/delete_many.go.tmpl"
			filed, // "crud/partials/delete_single.go.tmpl"
			filee, // "crud/partials/get.go.tmpl"
			filef, // "crud/partials/insert.go.tmpl"
			fileg, // "crud/partials/list.go.tmpl"
			fileh, // "crud/partials/loadrelated_manymany.go.tmpl"
			filei, // "crud/partials/loadrelated_manyone.go.tmpl"
			filej, // "crud/partials/loadrelated_onemany.go.tmpl"
			filek, // "crud/partials/merge.go.tmpl"
			filel, // "crud/partials/save.go.tmpl"
			filem, // "crud/partials/saverelated.go.tmpl"
			filen, // "crud/partials/update.go.tmpl"

		},
	}
	dirt := &embedded.EmbeddedDir{
		Filename:   "util",
		DirModTime: time.Unix(1527792315, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			fileu, // "util/util.go.tmpl"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "application"
		dir7, // "crud"
		dirt, // "util"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}
	dir7.ChildDirs = []*embedded.EmbeddedDir{
		dirb, // "crud/partials"

	}
	dirb.ChildDirs = []*embedded.EmbeddedDir{}
	dirt.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1527792315, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":              dir1,
			"application":   dir2,
			"crud":          dir7,
			"crud/partials": dirb,
			"util":          dirt,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"application/gen-service.sh.tmpl":            file3,
			"application/main.go.tmpl":                   file4,
			"bootstrap.go.tmpl":                          file5,
			"bootstrap_env.tmpl":                         file6,
			"crud/crud.go.tmpl":                          file8,
			"crud/hooks.go.tmpl":                         file9,
			"crud/models.go.tmpl":                        filea,
			"crud/partials/delete_many.go.tmpl":          filec,
			"crud/partials/delete_single.go.tmpl":        filed,
			"crud/partials/get.go.tmpl":                  filee,
			"crud/partials/insert.go.tmpl":               filef,
			"crud/partials/list.go.tmpl":                 fileg,
			"crud/partials/loadrelated_manymany.go.tmpl": fileh,
			"crud/partials/loadrelated_manyone.go.tmpl":  filei,
			"crud/partials/loadrelated_onemany.go.tmpl":  filej,
			"crud/partials/merge.go.tmpl":                filek,
			"crud/partials/save.go.tmpl":                 filel,
			"crud/partials/saverelated.go.tmpl":          filem,
			"crud/partials/update.go.tmpl":               filen,
			"crud/protobuf.proto.tmpl":                   fileo,
			"http.go.tmpl":                               filep,
			"rest.go.tmpl":                               fileq,
			"rest_hooks.go.tmpl":                         filer,
			"schema.sql.tmpl":                            files,
			"util/util.go.tmpl":                          fileu,
			"vuetify_actions.js.tmpl":                    filev,
			"vuetify_edit.vue.tmpl":                      filew,
			"vuetify_editor-field-checkbox.vue.tmpl":     filex,
			"vuetify_editor-field-date.vue.tmpl":         filey,
			"vuetify_editor-field-number.vue.tmpl":       filez,
			"vuetify_editor-field-password.vue.tmpl":     file10,
			"vuetify_editor-field-select-rel.vue.tmpl":   file11,
			"vuetify_editor-field-select.vue.tmpl":       file12,
			"vuetify_editor-field-textarea.vue.tmpl":     file13,
			"vuetify_editor-field-textfield.vue.tmpl":    file14,
			"vuetify_editor-field-time.vue.tmpl":         file15,
			"vuetify_editor-field-toggle.vue.tmpl":       file16,
			"vuetify_getters.js.tmpl":                    file17,
			"vuetify_index.js.tmpl":                      file18,
			"vuetify_list.vue.tmpl":                      file19,
			"vuetify_mutations.js.tmpl":                  file1a,
			"vuetify_routes.js.tmpl":                     file1b,
			"vuetify_types.js.tmpl":                      file1c,
		},
	})
}
