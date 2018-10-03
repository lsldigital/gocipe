package util

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
		Type: WidgetTypeTextField,
	},
}
