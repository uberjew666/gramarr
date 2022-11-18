package lidarr

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
)

var (
	apiRgx = regexp.MustCompile(`[a-z0-9]{32}`)
)

func NewClient(c Config) (*Client, error) {
	if c.Hostname == "" {
		return nil, fmt.Errorf("hostname is empty")
	}

	if match := apiRgx.MatchString(c.APIKey); !match {
		return nil, fmt.Errorf("api key is invalid format: %s", c.APIKey)
	}

	baseURL := createApiURL(c)

	r := resty.New()
	r.SetHostURL(baseURL)
	r.SetHeader("Accept", "application/json")
	r.SetQueryParam("apikey", c.APIKey)
	if c.Username != "" && c.Password != "" {
		r.SetBasicAuth(c.Username, c.Password)
	}

	client := Client{
		apiKey:     c.APIKey,
		maxResults: c.MaxResults,
		username:   c.Username,
		password:   c.Password,
		baseURL:    baseURL,
		client:     r,
	}
	return &client, nil
}

func createApiURL(c Config) string {
	c.Hostname = strings.TrimPrefix(c.Hostname, "http://")
	c.Hostname = strings.TrimPrefix(c.Hostname, "https://")
	c.URLBase = strings.TrimPrefix(c.URLBase, "/")

	u := url.URL{}
	if c.SSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if c.Port == 80 {
		u.Host = c.Hostname
	} else {
		u.Host = fmt.Sprintf("%s:%d", c.Hostname, c.Port)
	}

	if c.URLBase != "" {
		u.Path = fmt.Sprintf("%s/api/v1", c.URLBase)
	} else {
		u.Path = "/api/v1"
	}

	fmt.Println("The URL for Lidarr is", u.String())

	return u.String()
}

type Client struct {
	apiKey     string
	username   string
	password   string
	baseURL    string
	maxResults int
	client     *resty.Client
}

func (c *Client) SearchArtists(term string) ([]Artist, error) {
	resp, err := c.client.R().SetQueryParam("term", term).SetResult([]Artist{}).Get("artist/lookup")
	if err != nil {
		return nil, err
	}

	artists := *resp.Result().(*[]Artist)
	if len(artists) > c.maxResults {
		artists = artists[:c.maxResults]
	}
	return artists, nil
}

func (c *Client) GetFolders() ([]Folder, error) {
	resp, err := c.client.R().SetResult([]Folder{}).Get("rootfolder")
	if err != nil {
		return nil, err
	}

	folders := *resp.Result().(*[]Folder)
	return folders, nil
}

func (c *Client) GetProfile(endpoint string) ([]Profile, error) {

	resp, err := c.client.R().SetResult([]Profile{}).Get(endpoint)
	if err != nil {
		return nil, err
	}
	profile := *resp.Result().(*[]Profile)

	return profile, nil
}

func (c *Client) AddArtist(m Artist, metadataProfile int, qualityProfile int, path string) (artist Artist, err error) {

	request := AddArtistRequest{
		ArtistName:        m.ArtistName,
		CleanName:         m.CleanName,
		Images:            m.Images,
		QualityProfileID:  qualityProfile,
		MetadataProfileID: metadataProfile,
		ForeignArtistID:   m.ForeignArtistID,
		RootFolderPath:    path,
		Monitored:         true,
		AddOptions:        AddArtistOptions{SearchForMissingAlbums: true},
	}

	resp, err := c.client.R().SetBody(request).SetResult(Artist{}).Post("artist")
	if err != nil {
		fmt.Println(err)
		return
	}

	artist = *resp.Result().(*Artist)
	return
}
