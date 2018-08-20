package util

const (
	// WidgetTypeCheckbox represents widget of type checkbox
	WidgetTypeCheckbox = "checkbox"

	// WidgetTypeDate represents widget of type date
	WidgetTypeDate = "date"

	// WidgetTypeNumber represents widget of type number
	WidgetTypeNumber = "number"

	// WidgetTypePassword represents widget of type password
	WidgetTypePassword = "password"

	// WidgetTypeSelect represents widget of type select
	WidgetTypeSelect = "select"

	// WidgetTypeSelectRel represents widget of type select-rel
	WidgetTypeSelectRel = "select-rel"

	// WidgetTypeTextArea represents widget of type textarea
	WidgetTypeTextArea = "textarea"

	// WidgetTypeTextField represents widget of type textfield
	WidgetTypeTextField = "textfield"

	// WidgetTypeTime represents widget of type time
	WidgetTypeTime = "time"

	// WidgetTypeToggle represents widget of type toggle
	WidgetTypeToggle = "toggle"

	// RelationshipTypeManyMany represents a relationship of type Many-Many
	RelationshipTypeManyMany = "many-many"

	// RelationshipTypeOneOne represents a relationship of type One-One
	RelationshipTypeOneOne = "one-one"

	// RelationshipTypeOneMany represents a relationship of type One-Many
	RelationshipTypeOneMany = "one-many"

	// RelationshipTypeManyOne represents a relationship of type Many-One
	RelationshipTypeManyOne = "many-one"

	// PrimaryKeySerial indicates primary key - autogenerated number
	PrimaryKeySerial = "serial"

	// PrimaryKeyUUID indicates primary key - autogenerated string
	PrimaryKeyUUID = "uuid"

	// PrimaryKeyInt indicates primary key - manually assigned number
	PrimaryKeyInt = "int"

	// PrimaryKeyString indicates primary key - manually assigned string
	PrimaryKeyString = "string"
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

	// Bread describes options for Browse, Read, Edit, Add & Delete service generation
	Bread BreadOpts `json:"bread"`

	// Rest describes options for Rest generation
	Rest RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation
	Vuetify VuetifyOpts `json:"vuetify"`

	// Entities lists entities to be generated
	Entities []Entity `json:"entities"`
}

// HTTPOpts represents options for http function generation
type HTTPOpts struct {
	// Generate indicates whether or not to generate http serve function
	Generate bool `json:"generate"`

	// Port represents default port to run application
	Port string `json:"port"`
}

// BootstrapOpts represents options for bootstrap function generation
type BootstrapOpts struct {
	// Generate indicates whether or not to generate bootstrap
	Generate bool `json:"generate"`

	// NoDB indicates that database connection should not be generated (default false)
	NoDB bool `json:"no_db"`

	// NoGRPCWeb indicates that grpcweb server should not be generated (default false)
	NoGRPCWeb bool `json:"no_grpc_web"`

	// NoGRPCWire indicates that grpc server should not be generated (default false)
	NoGRPCWire bool `json:"no_grpc_wire"`

	// NoVersion indicates that version code should not be generated (default false)
	NoVersion bool `json:"no_version"`

	// Settings represent list of settings to load during bootstrap into main package
	Settings []BootstrapSetting `json:"settings"`

	// Assets indicates that we want to have an assets folder (using rice)
	Assets bool `json:"assets"`

	// HTTPPort represents port to listen to by default
	HTTPPort string `json:"http_port"`

	// GRPCPort represents port to grpc listen to by default
	GRPCPort string `json:"grpc_port"`
}

// BootstrapSetting represents a setting required by the application and loaded during bootstrap
type BootstrapSetting struct {
	// Name represents name of setting
	Name string `json:"name"`

	// Description gives information on the setting (useful to display errors if not found)
	Description string `json:"description"`

	// Public indicates if setting should be accessible from all packages
	Public bool `json:"public"`
}

// SchemaOpts represents options for schema generation
type SchemaOpts struct {
	// Create whether or not to generate CREATE TABLE
	Create bool `json:"create"`

	// Drop whether or not to generate DROP IF EXISTS before CREATE
	Drop bool `json:"drop"`

	// Aggregate whether or not to generate schema into single file instead of separate files
	Aggregate bool `json:"aggregate"`

	// Path indicates in which path to generate the schema sql file
	Path string `json:"path"`
}

// CrudOpts indicateds if crud functions should be generated
type CrudOpts struct {

	// Enable indicates if crud should be generated
	Generate bool `json:"generate"`

	// Hooks describes hooks options for CRUD generation
	Hooks CrudHooks `json:"hooks"`
}

// CrudHooks represents which crud hooks should be generated
type CrudHooks struct {

	// PreSave allows hook function to be executed before Save operation is performed
	PreSave bool `json:"pre_save"`

	// PostSave allows hook function to be executed after Save operation is performed
	PostSave bool `json:"post_save"`

	// PreRead allows hook function to be executed before Read operation is performed
	PreRead bool `json:"pre_read"`

	// PostRead allows hook function to be executed after Read operation is performed
	PostRead bool `json:"post_read"`

	// PreList allows hook function to be executed before List operation is performed
	PreList bool `json:"pre_list"`

	// PostList allows hook function to be executed after List operation is performed
	PostList bool `json:"post_list"`

	// PreDeleteSingle allows hook function to be executed before DeleteSingle operation is performed
	PreDeleteSingle bool `json:"pre_delete_single"`

	// PostDeleteSingle allows hook function to be executed after DeleteSingle operation is performed
	PostDeleteSingle bool `json:"post_delete_single"`

	// PreDeleteMany allows hook function to be executed before DeleteMany operation is performed
	PreDeleteMany bool `json:"pre_delete_many"`

	// PostDeleteMany allows hook function to be executed after DeleteMany operation is performed
	PostDeleteMany bool `json:"post_delete_many"`
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

	// Prefix indicates which prefix to use for routes
	Prefix string `json:"prefix"`

	// Hooks describes hooks options for REST generation
	Hooks ResourceHooks `json:"hooks"`
}

// BreadOpts represents which BREAD functions should be generated
type BreadOpts struct {
	// Generate indicates whether or not to generate the BREAD service
	Generate bool `json:"generate"`

	// Create indicates if code Add component of BREAD service, method Create, should be automatically generated
	Create bool `json:"create"`

	// Read indicates if code Read component of BREAD service, method Read, should be automatically generated
	Read bool `json:"read"`

	// List indicates if code Browse component of BREAD service, method List, should be automatically generated
	List bool `json:"list"`

	// Update indicates if code Edit component of BREAD service, method Update, should be automatically generated
	Update bool `json:"update"`

	// Delete indicates if code Delete component of BREAD service, method Delete, should be automatically generated
	Delete bool `json:"delete"`

	// Hooks describes hooks options for BREAD generation
	Hooks ResourceHooks `json:"hooks"`
}

// ResourceHooks represents which rest hooks should be generated
type ResourceHooks struct {

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

	// App represents the app for which the gocipe module will be generated
	App string `json:"app"`
}

// VuetifyEntityOpts represents per entity options for the vuetify generator
type VuetifyEntityOpts struct {
	// NoGenerate represents whether or not to generate vuetify assets
	NoGenerate bool `json:"no_generate"`

	// NotInMenu indicates whether or not to show entity in menu
	NotInMenu bool `json:"not_in_menu"`

	// Icon
	Icon string `json:"icon"`
}

// Entity represents a single entity to be generated
type Entity struct {
	// Name is the name of the entity
	Name string `json:"name"`

	// PrimaryKey indicates the nature of the primary key: serial (auto incremented number), uuid (auto generated string), int or string
	PrimaryKey string `json:"primary_key"`

	// Table is the name of the database table for the entity
	Table string `json:"table"`

	// TableConstraints represents an array of table constraints for the table definition
	TableConstraints []string `json:"table_constraints"`

	// Description is a description of the entity
	Description string `json:"description"`

	// Fields is a list of fields for the entity
	Fields []Field `json:"fields"`

	// Relationships represent relationship information between this entity and others
	Relationships []Relationship `json:"relationships"`

	// Schema describes options for Schema generation - overrides recipe level Schema config
	Schema *SchemaOpts `json:"schema"`

	// Crud describes options for CRUD generation - overrides recipe level Crud config
	CrudHooks *CrudHooks `json:"crud"`

	// Bread describes options for Bread generation - overrides recipe level Bread config
	Bread *BreadOpts `json:"bread"`

	// Rest describes options for Rest generation - overrides recipe level Rest config
	Rest *RestOpts `json:"rest"`

	// Vuetify describes options for Vuetify generation - overrides recipe level Vuetify config
	Vuetify VuetifyEntityOpts `json:"vuetify"`

	// DefaultSort is a sort string used while generating List() method in CRUD
	DefaultSort string `json:"default_sort"`
}

// Field describes a field contained in an entity
type Field struct {
	// Label is the label for the field
	Label string `json:"label"`

	// Property represents code property information for the field
	Property FieldProperty `json:"property"`

	// Schema represents schema information for the field
	Schema FieldSchema `json:"schema"`

	// EditWidget represents widget information for the field
	EditWidget EditWidgetOpts `json:"edit_widget"`

	// ListWidget represents widget information for the field
	ListWidget ListWidgetOpts `json:"list_widget"`

	// Filterable indicates if queries can be made using this field
	Filterable bool `json:"filterable"`
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

	// Default provides the default value for this field in database
	Default string `json:"default"`
}

// Relationship represents a relationship between this entity and another
type Relationship struct {
	// Entity is the name of the related entity
	Entity string `json:"entity"`

	// Type represents the type of relationship
	Type string `json:"type"`

	// Name represents the property name to be used for this relationship
	Name string `json:"name"`

	// JoinTable represents the other table in a many-many relationship
	JoinTable string `json:"join_table"`

	// ThisID represents the field in this entity used for the relationship
	ThisID string `json:"thisid"`

	// ThatID represents the field in the other entity used for the relationship
	ThatID string `json:"thatid"`

	// Eager indicates whether or not to eager load this relationship
	Eager bool `json:"eager"`
}

// EditWidgetOpts represents a UI widget for edit forms
type EditWidgetOpts struct {
	// Type indicates which widget type is represented
	Type string `json:"type"`

	// Options represents options listed by this widget
	Options []EditWidgetOption `json:"options"`

	// Target represents a target endpoint to pull data for this widget
	Target EditWidgetTarget `json:"target"`

	// Multiple indicates that the field accepts multiple values
	Multiple bool `json:"multiple"`
}

// EditWidgetOption represents an option for SelectRel widget type
type EditWidgetOption struct {
	// Value represents the stored value of the option
	Value string `json:"value"`
	// Label represents the displayed of the option
	Label string `json:"label"`
}

// EditWidgetTarget represents a target endpoint to pull data for this widget
type EditWidgetTarget struct {
	// Endpoint represents an endpoint to pull data from
	Endpoint string

	// Label which field to use for label on data endpoint
	Label string
}

// ListWidgetOpts represents a UI widget for listing tables
type ListWidgetOpts struct {
	// NoShowInList indicates whether or not to show field in listing
	Hide bool `json:"hide"`

	// Type indicates which widget type is represented
	Type string `json:"type"`
}
