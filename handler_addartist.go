package main

import (
	"fmt"
	"strings"

	"github.com/alcmoraes/gramarr/lidarr"

	"path/filepath"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleAddArtist(m *tb.Message) {
	e.CM.StartConversation(NewArtistShowConversation(e), m)
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
		c.ArtistQuery = m.Text

		Artists, err := c.env.lidarr.SearchArtists(c.ArtistQuery)
		c.ArtistResults = Artists

		// Search Service Failed
		if err != nil {
			SendError(c.env.Bot, m.Sender, "Failed to search Artist.")
			c.env.CM.StopConversation(c)
			return
		}

		// No Results
		if len(Artists) == 0 {
			msg := fmt.Sprintf("No Artists found with the title '%s'", EscapeMarkdown(c.ArtistQuery))
			Send(c.env.Bot, m.Sender, msg)
			c.env.CM.StopConversation(c)
			return
		}

		// Found some Artists! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d Artists:*", len(Artists)))
		for i, Art := range Artists {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, EscapeMarkdown(Art.String())))
		}
		Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickArtist(m)
	}
}

func (c *AddArtistConversation) AskPickArtist(m *tb.Message) Handler {

	// Send custom reply keyboard
	var options []string
	for _, Artist := range c.ArtistResults {
		options = append(options, fmt.Sprintf("%s", Artist))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which one would you like to download?", options)

	return func(m *tb.Message) {

		// Set the selected Artist
		for i := range options {
			if m.Text == options[i] {
				c.selectedArtist = &c.ArtistResults[i]
				break
			}
		}

		// Not a valid Artist selection
		if c.selectedArtist == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickArtist(m)
			return
		}

		c.currentStep = c.AskPickArtistSeason(m)
	}
}

func (c *AddTVShowConversation) AskPickTVShowSeason(m *tb.Message) Handler {

	// Send custom reply keyboard
	var options []string
	if len(c.selectedTVShowSeasons) > 0 {
		options = append(options, "Nope. I'm done!")
	}
	for _, Season := range c.selectedTVShow.Seasons {
		if len(c.selectedTVShowSeasons) > 0 {
			show := true
			for _, TVShowSeason := range c.selectedTVShowSeasons {
				if TVShowSeason.SeasonNumber == Season.SeasonNumber {
					show = false
				}
			}
			if show {
				options = append(options, fmt.Sprintf("%v", Season.SeasonNumber))
			}
		} else {
			options = append(options, fmt.Sprintf("%v", Season.SeasonNumber))
		}
	}
	options = append(options, "/cancel")
	if len(c.selectedTVShowSeasons) > 0 {
		SendKeyboardList(c.env.Bot, m.Sender, "Any other season?", options)
	} else {
		SendKeyboardList(c.env.Bot, m.Sender, "Which season would you like to download?", options)
	}

	return func(m *tb.Message) {

		if m.Text == "Nope. I'm done!" {
			for _, selectedTVShowSeason := range c.selectedTVShow.Seasons {
				selectedTVShowSeason.Monitored = false
				for _, chosenSeason := range c.selectedTVShowSeasons {
					if chosenSeason.SeasonNumber == selectedTVShowSeason.SeasonNumber {
						selectedTVShowSeason.Monitored = true
					}
				}
			}
			c.currentStep = c.AskPickTVShowQuality(m)
			return
		}

		// Set the selected Artist
		for i := range options {
			if m.Text == options[i] {
				c.selectedTVShowSeasons = append(c.selectedTVShowSeasons, *c.selectedTVShow.Seasons[i])
				break
			}
		}

		// Not a valid Artist selection
		if c.selectedTVShowSeasons == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickTVShowSeason(m)
			return
		}

		c.currentStep = c.AskPickTVShowSeason(m)
	}
}

func (c *AddTVShowConversation) AskPickTVShowQuality(m *tb.Message) Handler {

	profiles, err := c.env.lidarr.GetProfile("qualityprofile")

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
			c.currentStep = c.AskPickTVShowQuality(m)
			return
		}

		c.currentStep = c.AskPickTVShowLanguage(m)
	}
}

func (c *AddTVShowConversation) AskPickTVShowLanguage(m *tb.Message) Handler {

	languages, err := c.env.lidarr.GetProfile("languageprofile")

	// GetProfile Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get language profiles.")
		c.env.CM.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, LanguageProfile := range languages {
		options = append(options, fmt.Sprintf("%v", LanguageProfile.Name))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which language shall I look for?", options)

	return func(m *tb.Message) {
		// Set the selected option
		for i, opt := range options {
			if m.Text == opt {
				c.selectedLanguageProfile = &languages[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedLanguageProfile == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickTVShowLanguage(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddTVShowConversation) AskFolder(m *tb.Message) Handler {

	folders, err := c.env.lidarr.GetFolders()
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
		if c.selectedTVShow == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddTVShow(m)
	}
}

func (c *AddArtistConversation) AddArtist(m *tb.Message) {
	_, err := c.env.lidarr.AddArtist(*c.selectedArtist, c.selectedLanguageProfile.ID, c.selectedQualityProfile.ID, c.selectedFolder.Path)

	// Failed to add TV
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to add Artist.")
		c.env.CM.StopConversation(c)
		return
	}

	if c.selectedArtist.PosterURL != "" {
		photo := &tb.Photo{File: tb.FromURL(c.selectedArtist.PosterURL)}
		c.env.Bot.Send(m.Sender, photo)
	}

	// Notify User
	Send(c.env.Bot, m.Sender, "Artist has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added Artist '%s'", DisplayName(m.Sender), EscapeMarkdown(c.selectedArtist.String()))
	SendAdmin(c.env.Bot, c.env.Users.Admins(), adminMsg)

	c.env.CM.StopConversation(c)
}
