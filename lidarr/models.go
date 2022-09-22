package lidarr

import "fmt"

type Artist struct {
	ArtistName      string        `json:"artistName"`
	CleanName       string        `json:"cleanName"`
	PosterURL       string        `json:"remotePoster"`
	ForeignArtistID string        `json:"foreignArtistId"`
	Images          []ArtistImage `json:"images"`
}

type ArtistImage struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
}

type Profile struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Folder struct {
	Path      string `json:"path"`
	FreeSpace int64  `json:"freeSpace"`
	ID        int    `json:"id"`
}

type AddArtistRequest struct {
	ArtistName        string           `json:"artistName"`
	CleanName         string           `json:"cleanName"`
	Images            []ArtistImage    `json:"images"`
	QualityProfileID  int              `json:"qualityProfileId"`
	MetadataProfileID int              `json:"metadataProfileId"`
	ForeignArtistID   string           `json:"foreignArtistId"`
	RootFolderPath    string           `json:"rootFolderPath"`
	Monitored         bool             `json:"monitored"`
	AddOptions        AddArtistOptions `json:"addOptions"`
}

type AddArtistOptions struct {
	SearchForMissingAlbums bool `json:"searchForMissingAlbums"`
}
