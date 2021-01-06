package main

import (
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	chat "google.golang.org/api/chat/v1"
)

// ChatWorker is an interface for the work that the
// card service would have to do.
type ChatWorker interface {
	CreateIssueCard(*jira.Issue) (*chat.Message, error)
}

// ChatService is a wrapper for *chat.Service
type ChatService struct {
	*chat.Service
}

// Authorize authorizes the service to be used throughout
// the application
func (cs *ChatService) Authorize() {
	cs.Service = GetAuthdChatClient()
}

// CreateIssueCard takes the returned issue data and formats it into
// a card to be sent to google chat
func (cs ChatService) CreateIssueCard(issue *jira.Issue) (*chat.Message, error) {

	card := new(chat.Card)

	card.Header = &chat.CardHeader{
		Title:      issue.Fields.Summary,
		Subtitle:   issue.Key,
		ImageUrl:   issue.Fields.Type.IconURL,
		ImageStyle: "Image",
	}

	section := new(chat.Section)
	section.Widgets = append(section.Widgets, &chat.WidgetMarkup{
		TextParagraph: &chat.TextParagraph{
			Text: fmt.Sprintf(`<b>Type<\b>: %s<br><b>Priority<\b>: %s<br><b>Status<\b>: %s<br><b>Labels<\b>: %s<br>`,
				issue.Fields.Type.Description,
				issue.Fields.Priority.Name,
				issue.Fields.Status.Description,
				fmt.Sprintf("  %s", strings.Join(issue.Fields.Labels, "<br>  ")),
			),
		},
	})

	card.Sections = append(card.Sections, section)
	message := &chat.Message{Cards: []*chat.Card{card}}

	return message, nil
}
