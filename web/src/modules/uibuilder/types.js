//@ts-check
const Recipe = {
  tooltip: "Recipe represents a recipe to generate a project",
  definition: [
    {
      tooltip:
        "Container indicates whether or not container should be generated",
      label: "Bootstrap",
      type: "BootstrapOpts",
      json: "bootstrap"
    },

    {
      tooltip: "HTTP indicates whether http server code should be generated",
      label: "HTTP",
      type: "HTTPOpts",
      json: "http"
    },

    {
      tooltip: "Schema describes options for Schema generation",
      label: "Schema",
      type: "SchemaOpts",
      json: "schema"
    },

    {
      tooltip: "Crud describes options for Crud generation",
      label: "Crud",
      type: "CrudOpts",
      json: "crud"
    },

    {
      tooltip: "Rest describes options for Rest generation",
      label: "Rest",
      type: "RestOpts",
      json: "rest"
    },

    {
      tooltip: "Vuetify describes options for Vuetify generation",
      label: "Vuetify",
      type: "VuetifyOpts",
      json: "vuetify"
    },

    {
      tooltip: "Entities lists entities to be generated",
      label: "Entities",
      type: Entity
    }
  ]
};

const HTTPOpts = {
  tooltip: "HTTPOpts represents options for http function generation",
  definition: [
    {
      tooltip:
        "Generate indicates whether or not to generate http serve function",
      label: "Generate",
      type: "bool",
      json: "generate"
    },

    {
      tooltip: "Prefix indicates which prefix to use for routes",
      label: "Prefix",
      type: "string",
      json: "prefix"
    },

    {
      tooltip: "Port represents default port to run application",
      label: "Port",
      type: "string",
      json: "port"
    }
  ]
};

const BootstrapOpts = {
  tooltip: "BootstrapOpts represents options for bootstrap function generation",
  definition: [
    {
      tooltip: "Generate indicates whether or not to generate bootstrap",
      label: "Generate",
      type: "bool",
      json: "generate"
    },

    {
      tooltip:
        "Settings represent list of settings to load during bootstrap into main package",
      label: "Settings",
      type: "BootstrapSetting",
      json: "settings"
    }
  ]
};

const BootstrapSetting = {
  tooltip:
    "BootstrapSetting represents a setting required by the application and loaded during bootstrap",
  definition: [
    {
      tooltip: "Name represents name of setting",
      label: "Name",
      type: "string",
      json: "name"
    },

    {
      tooltip: "Type represents data type of setting",
      label: "Type",
      type: "string",
      json: "type"
    },

    {
      tooltip:
        "Description gives information on the setting (useful to display errors if not found)",
      label: "Description",
      type: "string",
      json: "description"
    },

    {
      tooltip:
        "FromDB indicates if setting comes from ENV variable or database (default)",
      label: "FromDB",
      type: "bool",
      json: "from_env"
    }
  ]
};

const SchemaOpts = {
  tooltip: "SchemaOpts represents options for schema generation",
  definition: [
    {
      tooltip: "Create whether or not to generate CREATE TABLE",
      label: "Create",
      type: "bool",
      json: "create"
    },

    {
      tooltip: "Drop whether or not to generate DROP IF EXISTS before CREATE",
      label: "Drop",
      type: "bool",
      json: "drop"
    },

    {
      tooltip:
        "Aggregate whether or not to generate schema into single file instead of separate files",
      label: "Aggregate",
      type: "bool",
      json: "aggregate"
    }
  ]
};

const CrudOpts = {
  tooltip: "CrudOpts represents which crud functions should be generated",
  definition: [
    {
      tooltip:
        "Create indicates whether or not function for INSERT should be generated",
      label: "Create",
      type: "bool",
      json: "create"
    },

    {
      tooltip:
        "Read indicates whether or not function for single select by id - SELECT WHERE id = id should be generated",
      label: "Read",
      type: "bool",
      json: "read"
    },

    {
      tooltip:
        "ReadList indicates whether or not function for list select should be generated",
      label: "ReadList",
      type: "bool",
      json: "read_list"
    },

    {
      tooltip:
        "Update indicates whether or not function for UPDATE should be generated",
      label: "Update",
      type: "bool",
      json: "update"
    },

    {
      tooltip:
        "Delete indicates whether or not function for DELETE should be generated",
      label: "Delete",
      type: "bool",
      json: "delete"
    },

    {
      tooltip: "Hooks describes hooks options for CRUD generation",
      label: "Hooks",
      type: "CrudHooks",
      json: "hooks"
    }
  ]
};

const CrudHooks = {
  tooltip: "CrudHooks represents which crud hooks should be generated",
  definition: [
    {
      tooltip:
        "PreCreate allows hook function to be executed before Create operation is performed",
      label: "PreCreate",
      type: "bool",
      json: "pre_create"
    },

    {
      tooltip:
        "PostCreate allows hook function to be executed after Create operation is performed",
      label: "PostCreate",
      type: "bool",
      json: "post_create"
    },

    {
      tooltip:
        "PreRead allows hook function to be executed before Read operation is performed",
      label: "PreRead",
      type: "bool",
      json: "pre_read"
    },

    {
      tooltip:
        "PostRead allows hook function to be executed after Read operation is performed",
      label: "PostRead",
      type: "bool",
      json: "post_read"
    },

    {
      tooltip:
        "PreList allows hook function to be executed before List operation is performed",
      label: "PreList",
      type: "bool",
      json: "pre_list"
    },

    {
      tooltip:
        "PostList allows hook function to be executed after List operation is performed",
      label: "PostList",
      type: "bool",
      json: "post_list"
    },

    {
      tooltip:
        "PreUpdate allows hook function to be executed before Update operation is performed",
      label: "PreUpdate",
      type: "bool",
      json: "pre_update"
    },

    {
      tooltip:
        "PostUpdate allows hook function to be executed after Update operation is performed",
      label: "PostUpdate",
      type: "bool",
      json: "post_update"
    },

    {
      tooltip:
        "PreDelete allows hook function to be executed before Delete operation is performed",
      label: "PreDelete",
      type: "bool",
      json: "pre_delete"
    },

    {
      tooltip:
        "PostDelete allows hook function to be executed after Delete operation is performed",
      label: "PostDelete",
      type: "bool",
      json: "post_delete"
    }
  ]
};

const RestOpts = {
  tooltip: "RestOpts represents which rest functions should be generated",
  definition: [
    {
      tooltip:
        "Create indicates if http endpoint for POST method should be generated",
      label: "Create",
      type: "bool",
      json: "create"
    },

    {
      tooltip:
        "Read indicates if http endpoint for GET method (by id for single entity) should be generated",
      label: "Read",
      type: "bool",
      json: "read"
    },

    {
      tooltip:
        "ReadList indicates if http endpoint for GET method (by filters for many entities) should be generated",
      label: "ReadList",
      type: "bool",
      json: "read_list"
    },

    {
      tooltip:
        "Update indicates if http endpoint for PUT method should be generated",
      label: "Update",
      type: "bool",
      json: "update"
    },

    {
      tooltip:
        "Delete indicates if http endpoint for DELETE method should be generated",
      label: "Delete",
      type: "bool",
      json: "delete"
    },

    {
      tooltip: "Hooks describes hooks options for REST generation",
      label: "Hooks",
      type: "RestHooks",
      json: "hooks"
    }
  ]
};

const RestHooks = {
  tooltip: "RestHooks represents which rest hooks should be generated",
  definition: [
    {
      tooltip:
        "PreCreate allows hook function to be executed before POST operations are done",
      label: "PreCreate",
      type: "bool",
      json: "pre_create"
    },

    {
      tooltip:
        "PostCreate allows hook function to be executed after POST operations are done",
      label: "PostCreate",
      type: "bool",
      json: "post_create"
    },

    {
      tooltip:
        "PreRead allows hook function to be executed before GET (single by id) operations are done",
      label: "PreRead",
      type: "bool",
      json: "pre_read"
    },

    {
      tooltip:
        "PostRead allows hook function to be executed after GET (single by id) operations are done",
      label: "PostRead",
      type: "bool",
      json: "post_read"
    },

    {
      tooltip:
        "PreList allows hook function to be executed before GET (many by filters) operations are done",
      label: "PreList",
      type: "bool",
      json: "pre_list"
    },

    {
      tooltip:
        "PostList allows hook function to be executed after GET (many by filters) operations are done",
      label: "PostList",
      type: "bool",
      json: "post_list"
    },

    {
      tooltip:
        "PreUpdate allows hook function to be executed before PUT operations are done",
      label: "PreUpdate",
      type: "bool",
      json: "pre_update"
    },

    {
      tooltip:
        "PostUpdate allows hook function to be executed after PUT operations are done",
      label: "PostUpdate",
      type: "bool",
      json: "post_update"
    },

    {
      tooltip:
        "PreDelete allows hook function to be executed before DELETE operations are done",
      label: "PreDelete",
      type: "bool",
      json: "pre_delete"
    },

    {
      tooltip:
        "PostDelete allows hook function to be executed after DELETE operations are done",
      label: "PostDelete",
      type: "bool",
      json: "post_delete"
    }
  ]
};

const VuetifyOpts = {
  tooltip: "VuetifyOpts represents options for vuetify generator",
  definition: [
    {
      tooltip: "Generate represents whether or not to generate vuetify assets",
      label: "Generate",
      type: "bool",
      json: "generate"
    },

    {
      tooltip: "Path represents the path where assets should be generated",
      label: "Path",
      type: "string",
      json: "path"
    }
  ]
};

const Entity = {
  tooltip: "Entity represents a single entity to be generated",
  definition: [
    {
      tooltip: "Name is the name of the entity",
      label: "Name",
      type: "string",
      json: "name"
    },

    {
      tooltip: "Table is the name of the database table for the entity",
      label: "Table",
      type: "string",
      json: "table"
    },

    {
      tooltip: "Description is a description of the entity",
      label: "Description",
      type: "string",
      json: "description"
    },

    {
      tooltip: "Fields is a list of fields for the entity",
      label: "Fields",
      type: "Field",
      json: "fields"
    },

    {
      tooltip:
        "Schema describes options for Schema generation - overrides recipe level Schema config",
      label: "Schema",
      type: "SchemaOpts",
      json: "schema"
    },

    {
      tooltip:
        "Crud describes options for Crud generation - overrides recipe level Crud config",
      label: "Crud",
      type: "CrudOpts",
      json: "crud"
    },

    {
      tooltip:
        "Rest describes options for Rest generation - overrides recipe level Rest config",
      label: "Rest",
      type: "RestOpts",
      json: "rest"
    },

    {
      tooltip:
        "Vuetify describes options for Vuetify generation - overrides recipe level Vuetify config",
      label: "Vuetify",
      type: "VuetifyOpts",
      json: "vuetify"
    }
  ]
};

const Field = {
  tooltip: "Field describes a field contained in an entity",
  definition: [
    {
      tooltip: "Label is the label for the field",
      label: "Label",
      type: "string",
      json: "label"
    },

    {
      tooltip:
        "Serialized is the name of the field for serialization (e.g. json)",
      label: "Serialized",
      type: "string",
      json: "serialized"
    },

    {
      tooltip: "Property represents code property information for the field",
      label: "Property",
      type: "FieldProperty",
      json: "property"
    },

    {
      tooltip: "Schema represents schema information for the field",
      label: "Schema",
      type: "FieldSchema",
      json: "schema"
    },

    {
      tooltip: "Relationship represents relationship information for the field",
      label: "Relationship",
      type: "FieldRelationship",
      json: "relationship"
    },

    {
      tooltip: "Widget represents widget information for the field",
      label: "Widget",
      type: "WidgetOptions",
      json: "widget"
    }
  ]
};

const FieldProperty = {
  tooltip: "FieldProperty represents code information for the field",
  definition: [
    {
      tooltip: "Name is the name of the property",
      label: "Name",
      type: "string",
      json: "name"
    },

    {
      tooltip: "Type is the data type of the property",
      label: "Type",
      type: "string",
      json: "type"
    }
  ]
};

const FieldSchema = {
  tooltip: "FieldSchema represents schema generation information for the field",
  definition: [
    {
      tooltip: "Field is the name of the field in database",
      label: "Field",
      type: "string",
      json: "field"
    },

    {
      tooltip: "Type is the data type for the field in database",
      label: "Type",
      type: "string",
      json: "type"
    },

    {
      tooltip:
        "Widget represents information to generate UI widgets - must be one of the widget types",
      label: "Widget",
      type: "interface"
    }
  ]
};

const FieldRelationship = {
  tooltip:
    "FieldRelationship represents a relationship between this entity and another",
  definition: [
    {
      tooltip: "Type represents the type of relationship",
      label: "Type",
      type: "string",
      json: "type"
    },
    {
      tooltip: "Target represents the target of the relationship",
      label: "Target",
      type: "FieldRelationshipTarget",
      json: "target"
    }
  ]
};

const FieldRelationshipTarget = {
  tooltip: "FieldRelationshipTarget represents a target for a relationship",
  definition: [
    {
      tooltip: "Entity represents the other entity in the relationship",
      label: "Entity",
      type: "string",
      json: "entity"
    },

    {
      tooltip: "Table represents the other table in the relationship",
      label: "Table",
      type: "string",
      json: "table"
    },

    {
      tooltip:
        "ThisID represents the field in this entity used for the relationship",
      label: "ThisID",
      type: "string",
      json: "thisid"
    },

    {
      tooltip:
        "ThatID represents the field in the other entity used for the relationship",
      label: "ThatID",
      type: "string",
      json: "thatid"
    }
  ]
};

const WidgetOpts = {
  tooltip: "WidgetOpts represents a UI widget",
  definition: [
    {
      tooltip: "Type indicates which widget type is represented",
      label: "Type",
      type: "string",
      json: "type"
    },

    {
      tooltip: "Options represents options listed by this widget",
      label: "Options",
      type: "WidgetOptions",
      json: "options"
    },

    {
      tooltip:
        "Target represents a target endpoint to pull data for this widget",
      label: "Target",
      type: "WidgetTarget",
      json: "target"
    }
  ]
};

const WidgetOptions = {
  tooltip: "WidgetOptions represents an option for SelectRel widget type",
  definition: [
    {
      tooltip: "Value represents the stored value of the option",
      label: "Value",
      type: "string",
      json: "value"
    },
    {
      tooltip: "Label represents the displayed of the option",
      label: "Label",
      type: "string",
      json: "label"
    }
  ]
};

const WidgetTarget = {
  tooltip:
    "WidgetTarget represents a target endpoint to pull data for this widget",
  definition: [
    {
      tooltip: "Endpoint represents an endpoint to pull data from",
      label: "Endpoint",
      type: "string"
    },

    {
      tooltip: "Label which field to use for label on data endpoint",
      label: "Label",
      type: "string"
    }
  ]
};

module.exports = {
  bootstrap: BootstrapOpts,
  http: HTTPOpts,
  schema: SchemaOpts,
  crud: CrudOpts,
  rest: RestOpts,
  vuetify: VuetifyOpts,
  settings: BootstrapSetting,
  hooks: CrudHooks,
  // hooks: RestHooks,
  fields: Field,
  // schema: SchemaOpts,
  // crud: CrudOpts,
  // rest: RestOpts,
  // vuetify: VuetifyOpts,
  property: FieldProperty,
  // schema: FieldSchema,
  relationship: FieldRelationship,
  widget: WidgetOptions,
  target: FieldRelationshipTarget,
  options: WidgetOptions,
  entities: Entity
  // target: WidgetTarget
};
