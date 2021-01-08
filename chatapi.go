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
	logger.Trace("Creating Issue Card")

	// Creating the card that will hold the whole message
	card := new(chat.Card)

	// Creating the card Header
	card.Header = &chat.CardHeader{
		Title:    issue.Fields.Summary,
		Subtitle: issue.Key,
	}

	// The peek data I'd like to include
	text := &chat.TextParagraph{
		Text: fmt.Sprintf(`<b>Assignee</b>: %s<br><b>Reporter</b>: %s<br><b>Type</b>: %s<br><b>Status</b>: %s<br>`,
			issue.Fields.Assignee.DisplayName,
			issue.Fields.Reporter.DisplayName,
			issue.Fields.Type.Name,
			issue.Fields.Status.Name,
		),
	}

	// for some reason our sprints are stored in a non standard way, so I have to do this
	// stuffs to filter out the one I want.
	sprints, err := issue.Fields.Unknowns.StringArray("customfield_10000")
	if err == nil {
		var name string

		for _, sprint := range sprints {
			sprintMeta := strings.Split(sprint, ",")
			if strings.Contains(sprintMeta[2], "ACTIVE") {
				name = strings.Split(sprintMeta[3], "=")[1]
			}
		}

		if name != "" {
			text.Text += fmt.Sprintf(`<b>Sprint</b>: %s<br>`, name)
		}
	}

	if len(issue.Fields.Labels) > 0 {
		text.Text += fmt.Sprintf("<b>Labels</b>:<br>%s<br>",
			fmt.Sprintf("  %s", strings.Join(issue.Fields.Labels, "<br>  ")),
		)
	}

	if len(issue.Fields.Components) > 0 {
		text.Text += fmt.Sprint("<b>Components</b>:<br>")

		for _, component := range issue.Fields.Components {
			text.Text += fmt.Sprintf("  %s<br>", component.Name)
		}
	}

	// Creating a section to put the peek data in
	textSection := new(chat.Section)
	textSection.Widgets = append(textSection.Widgets, &chat.WidgetMarkup{
		TextParagraph: text,
	})

	// button that sends you to the jira
	btns := []*chat.Button{{
		TextButton: &chat.TextButton{
			Text: "To Jira",
			OnClick: &chat.OnClick{
				OpenLink: &chat.OpenLink{
					Url: jiraBaseURL + "/browse/" + issue.Key,
				},
			},
		},
	}}

	// Creation a section to but the goto button
	btnSection := new(chat.Section)
	btnSection.Widgets = append(btnSection.Widgets, &chat.WidgetMarkup{
		Buttons: btns,
	})

	// adding the peek data section to the card
	card.Sections = append(card.Sections, textSection)
	card.Sections = append(card.Sections, btnSection)

	message := &chat.Message{Cards: []*chat.Card{card}}

	return message, nil
}
