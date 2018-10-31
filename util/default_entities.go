package util

var card = Entity{
	Name:        "Card",
	PrimaryKey:  "uuid",
	Description: "contains a list of cards to be used with decks",
	Vuetify: VuetifyEntityOpts{
		Icon: "pages",
	},
	Fields: []Field{
		Field{
			Label: "Deck Machine Name",
			Name:  "DeckMachineName",
			Type:  "string",
			EditWidget: EditWidgetOpts{
				Type: WidgetTypeSelect,
			},
		},
		Field{
			Label: "Position",
			Name:  "Position",
			Type:  "string",
			EditWidget: EditWidgetOpts{
				Type: WidgetTypeTextField,
			},
		},
	},
	Relationships: []Relationship{
		Relationship{
			Name:   "CardSchedule",
			Entity: "CardSchedule",
			Type:   RelationshipTypeManyOne,
		},
	},
	References: []Reference{
		Reference{
			Name: "Entity",
		},
	},
}

var cardSchedule = Entity{
	Name:        "CardSchedule",
	PrimaryKey:  "uuid",
	Description: "contains the schedule for the cards",
	Vuetify: VuetifyEntityOpts{
		Icon: "pages",
	},
	Fields: []Field{
		Field{
			Label: "Date Time",
			Name:  "DateTime",
			Type:  "time",
			EditWidget: EditWidgetOpts{
				Type: WidgetTypeTime,
			},
		},
		Field{
			Label: "Action",
			Name:  "Action",
			Type:  "string",
			EditWidget: EditWidgetOpts{
				Type: WidgetTypeTextField,
			},
		},
	},
	Relationships: []Relationship{
		Relationship{
			Name:   "Card",
			Entity: "Card",
			Type:   RelationshipTypeOneMany,
		},
	},
}
