package generators

const (
	// WidgetTypePassword indicates a widget of type Password
	WidgetTypePassword = "password"

	// WidgetTypeSelectRel indicates a widget of type Select-Relational
	WidgetTypeSelectRel = "select-rel"

	// WidgetTypeSelect indicates a widget of type Select
	WidgetTypeSelect = "select"

	// WidgetTypeTextarea indicates a widget of type Textarea
	WidgetTypeTextarea = "textarea"

	// WidgetTypeTextfield indicates a widget of type Textfield
	WidgetTypeTextfield = "textfield"
)

// Recipe represents a recipe to generate a project
type Recipe struct {
	// Container indicates whether or not container should be generated
	Bootstrap BootstrapOpts `json:"bootstrap"`

	// HTTP indicates whether http server code should be generated
	HTTP HTTPOpts `json:"http"`

	// Schema describes options for Schema generation
	Schema SchemaOpts `json:"schema"`

	// Crud describes options for Crud generation
	Crud CrudOpts `json:"crud"`

	// Rest describes options for Rest generation
	Rest RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation
	Vuetify VuetifyOpts `json:"vuetify"`

	// Entities lists entities to be generated
	Entities []Entity
}

// HTTPOpts represents options for http function generation
type HTTPOpts struct {
	// Generate indicates whether or not to generate http serve function
	Generate bool `json:"generate"`

	// Prefix indicates which prefix to use for routes
	Prefix string `json:"prefix"`

	// Port represents default port to run application
	Port string `json:"port"`
}

// BootstrapOpts represents options for bootstrap function generation
type BootstrapOpts struct {
	// Generate indicates whether or not to generate bootstrap
	Generate bool `json:"generate"`

	// Settings represent list of settings to load during bootstrap into main package
	Settings []BootstrapSetting `json:"settings"`
}

// BootstrapSetting represents a setting required by the application and loaded during bootstrap
type BootstrapSetting struct {
	// Name represents name of setting
	Name string `json:"name"`

	// Type represents data type of setting
	Type string `json:"type"`

	// Description gives information on the setting (useful to display errors if not found)
	Description string `json:"description"`

	// FromDB indicates if setting comes from ENV variable or database (default)
	FromDB bool `json:"from_env"`
}

// SchemaOpts represents options for schema generation
type SchemaOpts struct {
	// Create whether or not to generate CREATE TABLE
	Create bool `json:"create"`

	// Drop whether or not to generate DROP IF EXISTS before CREATE
	Drop bool `json:"drop"`

	// Aggregate whether or not to generate schema into single file instead of separate files
	Aggregate bool `json:"aggregate"`
}

// CrudOpts represents which crud functions should be generated
type CrudOpts struct {
	// Create indicates whether or not function for INSERT should be generated
	Create bool `json:"create"`

	// Read indicates whether or not function for single select by id - SELECT WHERE id = id should be generated
	Read bool `json:"read"`

	// ReadList indicates whether or not function for list select should be generated
	ReadList bool `json:"read_list"`

	// Update indicates whether or not function for UPDATE should be generated
	Update bool `json:"update"`

	// Delete indicates whether or not function for DELETE should be generated
	Delete bool `json:"delete"`

	// Hooks describes hooks options for CRUD generation
	Hooks CrudHooks `json:"hooks"`
}

// CrudHooks represents which crud hooks should be generated
type CrudHooks struct {

	// PreCreate allows hook function to be executed before Create operation is performed
	PreCreate bool `json:"pre_create"`

	// PostCreate allows hook function to be executed after Create operation is performed
	PostCreate bool `json:"post_create"`

	// PreRead allows hook function to be executed before Read operation is performed
	PreRead bool `json:"pre_read"`

	// PostRead allows hook function to be executed after Read operation is performed
	PostRead bool `json:"post_read"`

	// PreList allows hook function to be executed before List operation is performed
	PreList bool `json:"pre_list"`

	// PostList allows hook function to be executed after List operation is performed
	PostList bool `json:"post_list"`

	// PreUpdate allows hook function to be executed before Update operation is performed
	PreUpdate bool `json:"pre_update"`

	// PostUpdate allows hook function to be executed after Update operation is performed
	PostUpdate bool `json:"post_update"`

	// PreDelete allows hook function to be executed before Delete operation is performed
	PreDelete bool `json:"pre_delete"`

	// PostDelete allows hook function to be executed after Delete operation is performed
	PostDelete bool `json:"post_delete"`
}

// RestOpts represents which rest functions should be generated
type RestOpts struct {

	// Create indicates if http endpoint for POST method should be generated
	Create bool `json:"create"`

	// Read indicates if http endpoint for GET method (by id for single entity) should be generated
	Read bool `json:"read"`

	// ReadList indicates if http endpoint for GET method (by filters for many entities) should be generated
	ReadList bool `json:"read_list"`

	// Update indicates if http endpoint for PUT method should be generated
	Update bool `json:"update"`

	// Delete indicates if http endpoint for DELETE method should be generated
	Delete bool `json:"delete"`

	// Hooks describes hooks options for REST generation
	Hooks RestHooks `json:"hooks"`
}

// RestHooks represents which rest hooks should be generated
type RestHooks struct {

	// PreCreate allows hook function to be executed before POST operations are done
	PreCreate bool `json:"pre_create"`

	// PostCreate allows hook function to be executed after POST operations are done
	PostCreate bool `json:"post_create"`

	// PreRead allows hook function to be executed before GET (single by id) operations are done
	PreRead bool `json:"pre_read"`

	// PostRead allows hook function to be executed after GET (single by id) operations are done
	PostRead bool `json:"post_read"`

	// PreList allows hook function to be executed before GET (many by filters) operations are done
	PreList bool `json:"pre_list"`

	// PostList allows hook function to be executed after GET (many by filters) operations are done
	PostList bool `json:"post_list"`

	// PreUpdate allows hook function to be executed before PUT operations are done
	PreUpdate bool `json:"pre_update"`

	// PostUpdate allows hook function to be executed after PUT operations are done
	PostUpdate bool `json:"post_update"`

	// PreDelete allows hook function to be executed before DELETE operations are done
	PreDelete bool `json:"pre_delete"`

	// PostDelete allows hook function to be executed after DELETE operations are done
	PostDelete bool `json:"post_delete"`
}

// VuetifyOpts represents options for vuetify generator
type VuetifyOpts struct {
	// Generate represents whether or not to generate vuetify assets
	Generate bool `json:"generate"`

	// Path represents the path where assets should be generated
	Path string `json:"path"`
}

// Entity represents a single entity to be generated
type Entity struct {
	// Name is the name of the entity
	Name string `json:"name"`

	// Table is the name of the database table for the entity
	Table string `json:"table"`

	// Description is a description of the entity
	Description string `json:"description"`

	// Fields is a list of fields for the entity
	Fields []Field `json:"fields"`

	// Schema describes options for Schema generation - overrides recipe level Schema config
	Schema SchemaOpts `json:"schema"`

	// Crud describes options for Crud generation - overrides recipe level Crud config
	Crud CrudOpts `json:"crud"`

	// Rest describes options for Rest generation - overrides recipe level Rest config
	Rest RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation - overrides recipe level Vuetify config
	Vuetify VuetifyOpts `json:"vuetify"`
}

// Field describes a field contained in an entity
type Field struct {
	// Label is the label for the field
	Label string `json:"label"`

	// Serialized is the name of the field for serialization (e.g. json)
	Serialized string `json:"serialized"`

	// Property represents code property information for the field
	Property FieldProperty `json:"property"`

	// Schema represents schema information for the field
	Schema FieldSchema `json:"schema"`

	// Relationship represents relationship information for the field
	Relationship FieldRelationship `json:"relationship"`

	// Widget represents widget information for the field
	Widget WidgetOptions `json:"widget"`
}

// FieldProperty represents code information for the field
type FieldProperty struct {
	// Name is the name of the property
	Name string `json:"name"`

	// Type is the data type of the property
	Type string `json:"type"`
}

// FieldSchema represents schema generation information for the field
type FieldSchema struct {
	// Field is the name of the field in database
	Field string `json:"field"`

	// Type is the data type for the field in database
	Type string `json:"type"`

	// Widget represents information to generate UI widgets - must be one of the widget types
	Widget interface{}
}

// FieldRelationship represents a relationship between this entity and another
type FieldRelationship struct {
	// Type represents the type of relationship
	Type string `json:"type"`
	// Target represents the target of the relationship
	Target FieldRelationshipTarget `json:"target"`
}

// FieldRelationshipTarget represents a target for a relationship
type FieldRelationshipTarget struct {
	// Entity represents the other entity in the relationship
	Entity string `json:"entity"`

	// Table represents the other table in the relationship
	Table string `json:"table"`

	// ThisID represents the field in this entity used for the relationship
	ThisID string `json:"thisid"`

	// ThatID represents the field in the other entity used for the relationship
	ThatID string `json:"thatid"`
}

// WidgetOpts represents a UI widget
type WidgetOpts struct {
	// Type indicates which widget type is represented
	Type string `json:"type"`

	// Options represents options listed by this widget
	Options []WidgetOptions `json:"options"`

	// Target represents a target endpoint to pull data for this widget
	Target WidgetTarget `json:"target"`
}

// WidgetOptions represents an option for SelectRel widget type
type WidgetOptions struct {
	// Value represents the stored value of the option
	Value string `json:"value"`
	// Label represents the displayed of the option
	Label string `json:"label"`
}

// WidgetTarget represents a target endpoint to pull data for this widget
type WidgetTarget struct {
	// Endpoint represents an endpoint to pull data from
	Endpoint string

	// Label which field to use for label on data endpoint
	Label string
}
