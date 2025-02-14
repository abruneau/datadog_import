package widgets

import (
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

func (pc *PanelConvertor) newTextDefinition() (datadogV1.WidgetDefinition, error) {
	var markdown string
	if pc.panel.Title != "" {
		markdown = fmt.Sprintf("# %s\n", pc.panel.Title)
	}
	var mode, content string

	mode = pc.panel.Mode
	if pc.panel.Mode == "" {
		mode = pc.panel.Options.Mode
	}

	content = pc.panel.Content
	if pc.panel.Content == "" {
		content = pc.panel.Options.Content
	}

	if mode == "markdown" {
		markdown = markdown + content
	} else if mode == "html" {
		converter := md.NewConverter("", true, nil)
		m, err := converter.ConvertString(content)
		if err != nil {
			markdown = markdown + content
		}
		markdown = markdown + m

	}
	// markdown = markdown + content
	def := datadogV1.NewNoteWidgetDefinition(markdown, datadogV1.NOTEWIDGETDEFINITIONTYPE_NOTE)
	def.SetShowTick(false)
	def.SetBackgroundColor("white")
	return datadogV1.NoteWidgetDefinitionAsWidgetDefinition(def), nil
}
