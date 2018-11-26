package util

var fieldID = Field{
	Label: "ID",
	Name:  "ID",
	Type:  "string",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: true,
	},
}

var fieldStatus = Field{
	Label: "Status",
	Name:  "Status",
	Type:  "string",
	EditWidget: EditWidgetOpts{
		Type: WidgetTypeStatus,
		Options: []EditWidgetOption{
			EditWidgetOption{Text: "Draft", Value: StatusDraft},
			EditWidgetOption{Text: "Saved", Value: StatusSaved},
			EditWidgetOption{Text: "Published", Value: StatusPublished},
			EditWidgetOption{Text: "Unpublished", Value: StatusUnpublished},
		},
	},
	ListWidget: ListWidgetOpts{
		Type: WidgetTypeSelect,
	},
}

var fieldUserID = Field{
	Label: "UserID",
	Name:  "UserID",
	Type:  "string",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: true,
	},
}

var fieldSlug = Field{
	Label: "Slug",
	Name:  "Slug",
	Type:  "string",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: true,
		Type: WidgetTypeTextField,
	},
}

var contentField = Field{
	Label: "Lardwaz",
	Name:  "Content",
	Type:  "string",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: true,
	},
}

var fieldCreatedAt = Field{
	Label: "Created",
	Name:  "CreatedAt",
	Type:  "time",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: false,
		Type: WidgetTypeTime,
	},
}

var fieldUpdatedAt = Field{
	Label: "Updated",
	Name:  "UpdatedAt",
	Type:  "time",
	EditWidget: EditWidgetOpts{
		Hide: true,
	},
	ListWidget: ListWidgetOpts{
		Hide: false,
		Type: WidgetTypeTime,
	},
}

func fieldReferenceMakeFields(name string) (Field, Field) {
	var idField = Field{
		Label: name + "ID",
		Name:  name,
		Type:  "string",
	}

	var typeField = Field{
		Label: name + "Type",
		Name:  name + "Type",
		Type:  "string",
		EditWidget: EditWidgetOpts{
			Type: WidgetTypeSelect,
		},
	}

	return idField, typeField
}
