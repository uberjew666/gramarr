package main

import (
	"fmt"
	"strings"

	"github.com/uberjew666/gramarr/lidarr"

	"path/filepath"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleAddArtist(m *tb.Message) {
	e.CM.StartConversation(NewAddArtistConversation(e), m)
}

func NewAddArtistConversation(e *Env) *AddArtistConversation {
	return &AddArtistConversation{env: e}
}

type AddArtistConversation struct {
	currentStep             Handler
	artistQuery             string
	artistResults           []lidarr.Artist
	folderResults           []lidarr.Folder
	selectedArtist          *lidarr.Artist
	selectedQualityProfile  *lidarr.Profile
	selectedMetadataProfile *lidarr.Profile
	selectedFolder          *lidarr.Folder
	env                     *Env
}

func (c *AddArtistConversation) Run(m *tb.Message) {
	c.currentStep = c.AskArtist(m)
}

func (c *AddArtistConversation) Name() string {
	return "addartist"
}

func (c *AddArtistConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *AddArtistConversation) AskArtist(m *tb.Message) Handler {
	Send(c.env.Bot, m.Sender, "What Artist do you want to search for?")

	return func(m *tb.Message) {
		c.artistQuery = m.Text

		artists, err := c.env.Lidarr.SearchArtists(c.artistQuery)
		c.artistResults = artists

		// Search Service Failed
		if err != nil {
			SendError(c.env.Bot, m.Sender, "Failed to search Artist.")
			c.env.CM.StopConversation(c)
			return
		}

		// No Results
		if len(artists) == 0 {
			msg := fmt.Sprintf("No Artists found with the title '%s'", EscapeMarkdown(c.artistQuery))
			Send(c.env.Bot, m.Sender, msg)
			c.env.CM.StopConversation(c)
			return
		}

		// Found some artists! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d artists:*", len(artists)))
		for i, artist := range artists {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, EscapeMarkdown(artist.String())))
		}
		Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickArtist(m)
	}
}

func (c *AddArtistConversation) AskPickArtist(m *tb.Message) Handler {

	// Send custom reply keyboard
	var options []string
	for _, artist := range c.artistResults {
		options = append(options, fmt.Sprintf("%s", artist))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which one would you like to download?", options)

	return func(m *tb.Message) {

		// Set the selected Artist
		for i := range options {
			if m.Text == options[i] {
				c.selectedArtist = &c.artistResults[i]
				break
			}
		}

		// Not a valid Artist selection
		if c.selectedArtist == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickArtist(m)
			return
		}

		c.currentStep = c.AskPickArtistQuality(m)
	}
}

func (c *AddArtistConversation) AskPickArtistQuality(m *tb.Message) Handler {

	profiles, err := c.env.Lidarr.GetProfile("qualityprofile")

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
			c.currentStep = c.AskPickArtistQuality(m)
			return
		}

		c.currentStep = c.AskPickArtistMetadata(m)
	}
}

func (c *AddArtistConversation) AskPickArtistMetadata(m *tb.Message) Handler {

	profiles, err := c.env.Lidarr.GetProfile("metadataprofile")

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
			c.currentStep = c.AskPickArtistMetadata(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddArtistConversation) AskFolder(m *tb.Message) Handler {

	folders, err := c.env.Lidarr.GetFolders()
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
		if c.selectedArtist == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddArtist(m)
	}
}

func (c *AddArtistConversation) AddArtist(m *tb.Message) {
	_, err := c.env.Lidarr.AddArtist(*c.selectedArtist, c.selectedMetadataProfile.ID, c.selectedQualityProfile.ID, c.selectedFolder.Path)

	// Failed to add TV
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to add Artist.")
		c.env.CM.StopConversation(c)
		return
	}

	// Notify User
	Send(c.env.Bot, m.Sender, "Artist has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added Artist '%s'", DisplayName(m.Sender), EscapeMarkdown(c.selectedArtist.String()))
	SendAdmin(c.env.Bot, c.env.Users.Admins(), adminMsg)

	c.env.CM.StopConversation(c)
}
