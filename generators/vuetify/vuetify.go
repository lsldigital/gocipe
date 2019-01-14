package vuetify

import (
	"path"
	"path/filepath"

	"github.com/jinzhu/inflection"
	"github.com/lsldigital/gocipe/output"
	"github.com/lsldigital/gocipe/util"
)

// Generate returns generated vuetify components
func Generate(out *output.Output, r *util.Recipe) {
	if !r.Vuetify.Generate {
		// work.Waitgroup.Done()
		return
	}

	var (
		forms        []string
		menuEntities []util.Entity
	)
	dstPath := path.Join(util.WorkingDir, "/web/", r.Vuetify.App, "/src/gocipe")

	for _, entity := range r.Entities {
		if entity.Vuetify.NoGenerate {
			continue
		}

		if !entity.Vuetify.NoGenerate {
			menuEntities = append(menuEntities, entity)
		}

		// go func(entity util.Entity) {
		data := struct {
			Entity util.Entity
		}{entity}

		filePath := path.Join(dstPath, "/forms/")
		fileName := inflection.Plural(entity.Name)

		out.GenerateAndOverwrite("Vuetify List", "vuetify/forms/list.vue.tmpl", filepath.Join(filePath, fileName+"List.vue"), output.WithHeader, data)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"List")

		out.GenerateAndOverwrite("Vuetify Edit", "vuetify/forms/edit.vue.tmpl", filepath.Join(filePath, fileName+"Edit.vue"), output.WithHeader, data)
		forms = append(forms, inflection.Plural(data.Entity.Name)+"Edit")
	}
	// routes
	out.GenerateAndOverwrite("GenerateVuetify Routes", "vuetify/js/routes.js.tmpl", filepath.Join(dstPath, "/routes.js"), output.WithHeader, struct {
		Entities      []util.Entity
		GenerateDecks bool
	}{menuEntities, r.Decks.Generate})

	// decks
	if r.Decks.Generate {
		out.GenerateAndOverwrite("GenerateVuetify Routes Decks", "vuetify/decks/routes.js.tmpl", filepath.Join(dstPath, "/decks/routes.js"), output.WithHeader, struct {
			Decks         []util.DeckOpts
			GenerateDecks bool
		}{r.Decks.Decks, r.Decks.Generate})

		groupedDecks := make(map[string][]util.DeckOpts, 0)
		for _, deck := range r.Decks.Decks {
			if deck.Vuetify.NoGenerate {
				continue
			}
			groupedDecks[deck.Group] = append(groupedDecks[deck.Group], deck)
			entities := make([]util.Entity, 0)
			for _, name := range deck.EntityTypeWhitelist {
				e, err := r.GetEntity(name)
				if err != nil {
					continue
				}
				entities = append(entities, *e)
			}
			out.GenerateAndOverwrite("GenerateVuetify Deck", "vuetify/decks/form.vue.tmpl", filepath.Join(dstPath, "/decks/"+deck.Name+".vue"), output.WithHeader, struct {
				Deck     util.DeckOpts
				Entities []util.Entity
			}{deck, entities})
		}

		out.GenerateAndOverwrite("GenerateVuetify Decks", "vuetify/decks/home.vue.tmpl", filepath.Join(dstPath, "/decks/home.vue"), output.WithHeader, struct {
			Decks map[string][]util.DeckOpts
		}{groupedDecks})
	}

	widgets := map[string]string{
		"EditWidgetIcon":       "edit/Icon.vue",
		"EditWidgetImagefield": "edit/Imagefield.vue",
		"EditWidgetMap":        "edit/Map.vue",
		"EditWidgetStatus":     "edit/Status.vue",
		"EditWidgetSelect":     "edit/Select.vue",
		"EditWidgetSelectRel":  "edit/SelectRel.vue",
		"EditWidgetTextarea":   "edit/Textarea.vue",
		"EditWidgetTextfield":  "edit/Textfield.vue",
		"EditWidgetTime":       "edit/Time.vue",
		"EditWidgetToggle":     "edit/Toggle.vue",
		"ListWidgetSelect":     "list/Select.vue",
		"ListWidgetTime":       "list/Time.vue",
		"ListWidgetToggle":     "list/Toggle.vue",
	}

	for _, file := range widgets {
		out.GenerateAndOverwrite("GenerateVuetify Widgets", filepath.Join("vuetify/widgets/", file+".tmpl"), filepath.Join(dstPath, "/widgets/", file), output.WithHeader, nil)
	}
	// components
	out.GenerateAndOverwrite("GenerateVuetify Registration", "vuetify/js/components-registration.js.tmpl", filepath.Join(dstPath, "/components-registration.js"), output.WithHeader, struct {
		Widgets map[string]string
		Forms   []string
	}{Widgets: widgets, Forms: forms})

	lardwaz := map[string]string{
		// "RendererProjectNameSCSS":   "lardwaz/renderer/ProjectName.scss",
		"README":                    "lardwaz/README.md",
		"StoreActions":              "lardwaz/store/actions.js",
		"StoreGetters":              "lardwaz/store/getters.js",
		"StoreIndex":                "lardwaz/store/index.js",
		"StoreMutations":            "lardwaz/store/mutations.js",
		"BlocksLardwaz":             "lardwaz/views/Lardwaz.vue",
		"BlocksFooterBlock":         "lardwaz/views/blocks/FooterBlock.vue",
		"BlocksGalleryRender":       "lardwaz/views/blocks/GalleryRender.vue",
		"BlocksHeadingBlock":        "lardwaz/views/blocks/HeadingBlock.vue",
		"BlocksImageRender":         "lardwaz/views/blocks/ImageRender.vue",
		"BlocksMarkdownBlock":       "lardwaz/views/blocks/MarkdownBlock.vue",
		"BlocksQuoteRender":         "lardwaz/views/blocks/QuoteRender.vue",
		"BlocksTextBlock":           "lardwaz/views/blocks/TextBlock.vue",
		"BlocksTextareaRender":      "lardwaz/views/blocks/TextareaRender.vue",
		"BlocksFooterRender":        "lardwaz/views/blocks/FooterRender.vue",
		"BlocksHTMLBlock":           "lardwaz/views/blocks/HTMLBlock.vue",
		"BlocksHeadingRender":       "lardwaz/views/blocks/HeadingRender.vue",
		"BlocksLegacyBlock":         "lardwaz/views/blocks/LegacyBlock.vue",
		"BlocksMarkdownRender":      "lardwaz/views/blocks/MarkdownRender.vue",
		"BlocksRelatedBlock":        "lardwaz/views/blocks/RelatedBlock.vue",
		"BlocksTextRender":          "lardwaz/views/blocks/TextRender.vue",
		"BlocksYoutubeBlock":        "lardwaz/views/blocks/YoutubeBlock.vue",
		"BlocksGalleryBlock":        "lardwaz/views/blocks/GalleryBlock.vue",
		"BlocksHTMLRender":          "lardwaz/views/blocks/HTMLRender.vue",
		"BlocksImageBlock":          "lardwaz/views/blocks/ImageBlock.vue",
		"BlocksLegacyRender":        "lardwaz/views/blocks/LegacyRender.vue",
		"BlocksQuoteBlock":          "lardwaz/views/blocks/QuoteBlock.vue",
		"BlocksRelatedRender":       "lardwaz/views/blocks/RelatedRender.vue",
		"BlocksTextareaBlock":       "lardwaz/views/blocks/TextareaBlock.vue",
		"BlocksYoutubeRender":       "lardwaz/views/blocks/YoutubeRender.vue",
		"BlockIndicator":            "lardwaz/views/blocks/IndicatorBlock.vue",
		"BlockIndicatorRender":      "lardwaz/views/blocks/IndicatorRender.vue",
		"BlockTranscript":           "lardwaz/views/blocks/TranscriptBlock.vue",
		"BlockTranscriptRender":     "lardwaz/views/blocks/TranscriptRender.vue",
		"ComponentsBlockEditor":     "lardwaz/views/components/BlockEditor.vue",
		"ComponentsBlockPicker":     "lardwaz/views/components/BlockPicker.vue",
		"ComponentsInformation":     "lardwaz/views/components/Information.vue",
		"MailAdBlock":               "lardwaz/views/blocks/mail/MailAdBlock.vue",
		"MailArticlesListingRender": "lardwaz/views/blocks/mail/MailArticlesListingRender.vue",
		"MailColumnRender":          "lardwaz/views/blocks/mail/MailColumnRender.vue",
		"MailHeaderRender":          "lardwaz/views/blocks/mail/MailHeaderRender.vue",
		"MailImageTextRender":       "lardwaz/views/blocks/mail/MailImageTextRender.vue",
		"MailJumbotronRender":       "lardwaz/views/blocks/mail/MailJumbotronRender.vue",
		"MailTextHeadingRender":     "lardwaz/views/blocks/mail/MailTextHeadingRender.vue",
		"MailAdRender":              "lardwaz/views/blocks/mail/MailAdRender.vue",
		"MailColumnBlock":           "lardwaz/views/blocks/mail/MailColumnBlock.vue",
		"MailHeaderBlock":           "lardwaz/views/blocks/mail/MailHeaderBlock.vue",
		"MailImageTextBlock":        "lardwaz/views/blocks/mail/MailImageTextBlock.vue",
		"MailJumbotronBlock":        "lardwaz/views/blocks/mail/MailJumbotronBlock.vue",
		"MailTextHeadingBlock":      "lardwaz/views/blocks/mail/MailTextHeadingBlock.vue",
		"MailArticlesListingBlock":  "lardwaz/views/blocks/mail/MailArticlesListingBlock.vue",
		"MailColumnMjml":            "lardwaz/views/blocks/mail/MailColumnMjml.vue",
		"MailHeaderMjml":            "lardwaz/views/blocks/mail/MailHeaderMjml.vue",
		"MailImageTextMjml":         "lardwaz/views/blocks/mail/MailImageTextMjml.vue",
		"MailJumbotronMjml":         "lardwaz/views/blocks/mail/MailJumbotronMjml.vue",
		"MailTextHeadingMjml":       "lardwaz/views/blocks/mail/MailTextHeadingMjml.vue",
		"ConfigPreviewComponent":    "lardwaz-config/Preview.vue",
		"ConfigPageComponents":      "lardwaz-config/page-components.js",
		"ConfigRenderer":            "lardwaz-config/Renderer.vue",
	}

	for name, file := range lardwaz {
		out.GenerateAndOverwrite("GenerateVuetify Lardwaz "+name, filepath.Join("vuetify/modules", file+".tmpl"), filepath.Join(util.WorkingDir, "web", r.Vuetify.App, "modules", file), output.WithHeader, struct {
			Recipe *util.Recipe
		}{r})
	}

	staticFiles := map[string]string{
		"shared-ui/AppFooter.vue":         "shared-ui/AppFooter.vue",
		"shared-ui/AppNavigation.vue":     "shared-ui/AppNavigation.vue",
		"shared-ui/AppToolbar.vue":        "shared-ui/AppToolbar.vue",
		"shared-ui/NotFound.vue":          "shared-ui/NotFound.vue",
		"shared-ui/PageHome.vue":          "shared-ui/PageHome.vue",
		"shared-ui/Authenticated.vue":     "shared-ui/Authenticated.vue",
		"shared-ui/Login.vue":             "shared-ui/Login.vue",
		"store/modules/auth/index.js":     "store/modules/auth/index.js",
		"store/modules/auth/getters.js":   "store/modules/auth/getters.js",
		"store/modules/auth/actions.js":   "store/modules/auth/actions.js",
		"store/modules/auth/mutations.js": "store/modules/auth/mutations.js",
		"store/modules/auth/types.js":     "store/modules/auth/types.js",
	}

	for src, target := range staticFiles {
		output.GenerateAndSave("Vuetify Static Files", "vuetify/"+src+".tmpl", filepath.Join(dstPath, target), output.WithHeader, nil)
	}

	output.GenerateAndSave("Vuetify Router", "vuetify/js/router.js.tmpl", filepath.Join(dstPath, "../router.js"), output.WithHeader, nil)
	output.GenerateAndSave("Vuetify Store", "vuetify/js/store.js.tmpl", filepath.Join(dstPath, "../store.js"), output.WithHeader, nil)
	output.GenerateAndSave("Vuetify App", "vuetify/shared-ui/App.vue.tmpl", filepath.Join(dstPath, "../App.vue"), output.WithHeader, nil)

	// output.GenerateAndSave("Vuetify", "vuetify/store/index.js.tmpl", path+"gocipe/store/index.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/actions.js.tmpl", path+"gocipe/store/actions.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/getters.js.tmpl", path+"gocipe/store/getters.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/mutations.js.tmpl", path+"gocipe/store/mutations.js", nil, false)
	// output.GenerateAndSave("Vuetify", "vuetify/store/types.js.tmpl", path+"gocipe/store/types.js", nil, false)

}
