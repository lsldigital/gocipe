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

const Settings = {
  tooltip: "Entity represents a single entity to be generated",
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

module.exports = {
  entities: Entity,
  settings: Settings
};
