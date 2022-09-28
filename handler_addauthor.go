package main

import (
	"fmt"
	"strings"

	"path/filepath"

	"github.com/uberjew666/gramarr/readarr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleAddAuthor(m *tb.Message) {
	e.CM.StartConversation(NewAddAuthorConversation(e), m)
}

func NewAddAuthorConversation(e *Env) *AddAuthorConversation {
	return &AddAuthorConversation{env: e}
}

type AddAuthorConversation struct {
	currentStep             Handler
	authorQuery             string
	authorResults           []readarr.AuthorResource
	folderResults           []readarr.Folder
	selectedAuthor          *readarr.AuthorResource
	selectedQualityProfile  *readarr.Profile
	selectedMetadataProfile *readarr.Profile
	selectedFolder          *readarr.Folder
	env                     *Env
}

func (c *AddAuthorConversation) Run(m *tb.Message) {
	c.currentStep = c.AskAuthor(m)
}

func (c *AddAuthorConversation) Name() string {
	return "addauthor"
}

func (c *AddAuthorConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *AddAuthorConversation) AskAuthor(m *tb.Message) Handler {
	Send(c.env.Bot, m.Sender, "What author do you want to search for?")

	return func(m *tb.Message) {
		c.authorQuery = m.Text

		authors, err := c.env.Readarr.SearchAuthors(c.authorQuery)
		c.authorResults = authors

		// Search Service Failed
		if err != nil {
			SendError(c.env.Bot, m.Sender, "Failed to search author.")
			c.env.CM.StopConversation(c)
			return
		}

		// No Results
		if len(authors) == 0 {
			msg := fmt.Sprintf("No authors found with the title '%s'", EscapeMarkdown(c.authorQuery))
			Send(c.env.Bot, m.Sender, msg)
			c.env.CM.StopConversation(c)
			return
		}

		// Found some authors! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d authors:*", len(authors)))
		for i, author := range authors {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, EscapeMarkdown(author.String())))
		}
		Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickAuthor(m)
	}
}

func (c *AddAuthorConversation) AskPickAuthor(m *tb.Message) Handler {

	// Send custom reply keyboard
	var options []string
	for _, author := range c.authorResults {
		options = append(options, fmt.Sprintf("%s", author))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which one would you like to download?", options)

	return func(m *tb.Message) {

		// Set the selected author
		for i := range options {
			if m.Text == options[i] {
				c.selectedAuthor = &c.authorResults[i]
				break
			}
		}

		// Not a valid author selection
		if c.selectedAuthor == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickAuthor(m)
			return
		}

		c.currentStep = c.AskPickAuthorQuality(m)
	}
}

func (c *AddAuthorConversation) AskPickAuthorQuality(m *tb.Message) Handler {

	profiles, err := c.env.Readarr.GetProfile("qualityprofile")

	// GetProfile Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get quality profiles.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, QualityProfile := range profiles {
		options = append(options, fmt.Sprintf("%v", QualityProfile.Name))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which quality shall I look for?", options)

	return func(m *tb.Message) {
		// Set the selected option
		for i := range options {
			if m.Text == options[i] {
				c.selectedQualityProfile = &profiles[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedQualityProfile == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickAuthorQuality(m)
			return
		}

		c.currentStep = c.AskPickAuthorMetadata(m)
	}
}

func (c *AddAuthorConversation) AskPickAuthorMetadata(m *tb.Message) Handler {

	profiles, err := c.env.Readarr.GetProfile("metadataprofile")

	// GetProfile Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get metadata profiles.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, MetadataProfile := range profiles {
		options = append(options, fmt.Sprintf("%v", MetadataProfile.Name))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which metadata shall I look for?", options)

	return func(m *tb.Message) {
		// Set the selected option
		for i, opt := range options {
			if m.Text == opt {
				c.selectedMetadataProfile = &profiles[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedMetadataProfile == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickAuthorMetadata(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddAuthorConversation) AskFolder(m *tb.Message) Handler {

	folders, err := c.env.Readarr.GetFolders()
	c.folderResults = folders

	// GetFolders Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get folders.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// No Results
	if len(folders) == 0 {
		SendError(c.env.Bot, m.Sender, "No destination folders found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Found folders!

	// Send the results
	var msg []string
	msg = append(msg, fmt.Sprintf("*Found %d folders:*", len(folders)))
	for i, folder := range folders {
		msg = append(msg, fmt.Sprintf("%d) %s", i+1, EscapeMarkdown(filepath.Base(folder.Path))))
	}
	Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))

	// Send the custom reply keyboard
	var options []string
	for _, folder := range folders {
		options = append(options, fmt.Sprintf("%s", filepath.Base(folder.Path)))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which folder should it download to?", options)

	return func(m *tb.Message) {
		// Set the selected folder
		for i, opt := range options {
			if m.Text == opt {
				c.selectedFolder = &c.folderResults[i]
				break
			}
		}

		// Not a valid folder selection
		if c.selectedAuthor == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddAuthor(m)
	}
}

func (c *AddAuthorConversation) AddAuthor(m *tb.Message) {
	_, err := c.env.Readarr.AddAuthor(*c.selectedAuthor, c.selectedMetadataProfile.ID, c.selectedQualityProfile.ID, c.selectedFolder.Path)

	// Failed to add author
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to add Author.")
		c.env.CM.StopConversation(c)
		return
	}

	// Notify User
	Send(c.env.Bot, m.Sender, "Author has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added Author '%s'", DisplayName(m.Sender), EscapeMarkdown(c.selectedAuthor.String()))
	SendAdmin(c.env.Bot, c.env.Users.Admins(), adminMsg)

	c.env.CM.StopConversation(c)
}
