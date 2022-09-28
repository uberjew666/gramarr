package readarr

import (
	"fmt"
	"time"
)

func (m AuthorResource) String() string {
	if m.Disambiguation != "" {
		return fmt.Sprintf("%s (%s)", m.AuthorName, m.Disambiguation)
	} else {
		return m.AuthorName
	}
}

type BookResource struct {
	ID             int                    `json:"id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	AuthorTitle    string                 `json:"authorTitle,omitempty"`
	SeriesTitle    string                 `json:"seriesTitle,omitempty"`
	Disambiguation string                 `json:"disambiguation,omitempty"`
	Overview       string                 `json:"overview,omitempty"`
	AuthorID       int                    `json:"authorId,omitempty"`
	ForeignBookID  string                 `json:"foreignBookId,omitempty"`
	TitleSlug      string                 `json:"titleSlug,omitempty"`
	Monitored      bool                   `json:"monitored,omitempty"`
	AnyEditionOk   bool                   `json:"anyEditionOk,omitempty"`
	Ratings        Ratings                `json:"ratings,omitempty"`
	ReleaseDate    time.Time              `json:"releaseDate,omitempty"`
	PageCount      int                    `json:"pageCount,omitempty"`
	Genres         []string               `json:"genres,omitempty"`
	Author         AuthorResource         `json:"author,omitempty"`
	Images         []MediaCover           `json:"images,omitempty"`
	Links          []Links                `json:"links,omitempty"`
	Statistics     BookStatisticsResource `json:"statistics,omitempty"`
	Added          time.Time              `json:"added,omitempty"`
	AddOptions     AddBookOptions         `json:"addOptions,omitempty"`
	RemoteCover    string                 `json:"remoteCover,omitempty"`
	Editions       []EditionResource      `json:"editions,omitempty"`
	Grabbed        bool                   `json:"grabbed,omitempty"`
}

type Book struct {
	ID               int            `json:"id,omitempty"`
	AuthorMetadataID int            `json:"authorMetadataId,omitempty"`
	ForeignBookID    string         `json:"foreignBookId,omitempty"`
	TitleSlug        string         `json:"titleSlug,omitempty"`
	Title            string         `json:"title,omitempty"`
	ReleaseDate      time.Time      `json:"releaseDate,omitempty"`
	Links            []Links        `json:"links,omitempty"`
	Genres           []string       `json:"genres,omitempty"`
	Ratings          Ratings        `json:"ratings,omitempty"`
	CleanName        string         `json:"cleanName,omitempty"`
	Monitored        bool           `json:"monitored,omitempty"`
	AnyEditionOk     bool           `json:"anyEditionOk,omitempty"`
	LastInfoSync     time.Time      `json:"lastInfoSync,omitempty"`
	Added            time.Time      `json:"added,omitempty"`
	AddOptions       AddBookOptions `json:"addOptions,omitempty"`
}

type Ratings struct {
	Value float64 `json:"value,omitempty"`
	Votes int     `json:"votes,omitempty"`
}

type MediaCover struct {
	CoverType string `json:"coverType,omitempty"`
	Extension string `json:"extension,omitempty"`
	URL       string `json:"url,omitempty"`
}

type Links struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type BookStatisticsResource struct {
	BookCount      int `json:"bookCount,omitempty"`
	BookFileCount  int `json:"bookFileCount,omitempty"`
	SizeOnDisk     int `json:"sizeOnDisk,omitempty"`
	TotalBookCount int `json:"totalBookCount,omitempty"`
}

type AuthorResource struct {
	ID                  int                     `json:"id,omitempty"`
	AuthorMetadataID    int                     `json:"authorMetadataId,omitempty"`
	Status              string                  `json:"status,omitempty"`
	Ended               bool                    `json:"ended,omitempty"`
	AuthorName          string                  `json:"authorName,omitempty"`
	AuthorNameLastFirst string                  `json:"authorNameLastFirst,omitempty"`
	ForeignAuthorID     string                  `json:"foreignAuthorId,omitempty"`
	TitleSlug           string                  `json:"titleSlug,omitempty"`
	Overview            string                  `json:"overview,omitempty"`
	Disambiguation      string                  `json:"disambiguation,omitempty"`
	Links               []Links                 `json:"links,omitempty"`
	NextBook            Book                    `json:"nextBook,omitempty"`
	LastBook            Book                    `json:"lastBook,omitempty"`
	Images              []MediaCover            `json:"images,omitempty"`
	RemoteCover         string                  `json:"remoteCover,omitempty"`
	Path                string                  `json:"path,omitempty"`
	QualityProfileID    int                     `json:"qualityProfileId,omitempty"`
	MetadataProfileID   int                     `json:"metadataProfileId,omitempty"`
	Monitored           bool                    `json:"monitored,omitempty"`
	RootFolderPath      string                  `json:"rootFolderPath,omitempty"`
	Genres              []string                `json:"genres,omitempty"`
	CleanName           string                  `json:"cleanName,omitempty"`
	SortName            string                  `json:"sortName,omitempty"`
	SortNameLastFirst   string                  `json:"sortNameLastFirst,omitempty"`
	Tags                []int                   `json:"tags,omitempty"`
	Added               time.Time               `json:"added,omitempty"`
	AddOptions          AddAuthorOptions        `json:"addOptions,omitempty"`
	Ratings             Ratings                 `json:"ratings,omitempty"`
	Statistics          AuthorStatisticResource `json:"statistics,omitempty"`
}

type EditionResource struct {
	ID             int          `json:"id,omitempty"`
	BookID         int          `json:"bookId,omitempty"`
	ForeignBookID  string       `json:"foreignBookId,omitempty"`
	TitleSlug      string       `json:"titleSlug,omitempty"`
	ISBN13         string       `json:"isbn13,omitempty"`
	ASIN           string       `json:"asin,omitempty"`
	Title          string       `json:"title,omitempty"`
	Language       string       `json:"language,omitempty"`
	Overview       string       `json:"overview,omitempty"`
	Format         string       `json:"format,omitempty"`
	IsEbook        bool         `json:"isEbook,omitempty"`
	Disambiguation string       `json:"disambiguation,omitempty"`
	Publisher      string       `json:"publisher,omitempty"`
	PageCount      int          `json:"pageCount,omitempty"`
	ReleaseDate    time.Time    `json:"releaseDate,omitempty"`
	Images         []MediaCover `json:"images,omitempty"`
	Links          []Links      `json:"links,omitempty"`
	Ratings        Ratings      `json:"ratings,omitempty"`
	Monitored      bool         `json:"monitored,omitempty"`
	ManualAdd      bool         `json:"manualAdd,omitempty"`
	RemoteCover    string       `json:"remoteCover,omitempty"`
	Grabbed        bool         `json:"grabbed,omitempty"`
}

type Profile struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type Folder struct {
	Path      string `json:"path,omitempty"`
	FreeSpace int64  `json:"freeSpace,omitempty"`
	ID        int    `json:"id,omitempty"`
}

type AddBookOptions struct {
	SearchForBook bool `json:"searchForBook"`
}

type AddAuthorOptions struct {
	Monitor               string   `json:"monitor,omitempty"`
	BooksToMonitor        []string `json:"booksToMonitor,omitempty"`
	Monitored             bool     `json:"monitored,omitempty"`
	SearchForMissingBooks bool     `json:"searchForMissingBooks,omitempty"`
}

type AuthorStatisticResource struct {
	BookFileCount      int `json:"bookFileCount,omitempty"`
	BookCount          int `json:"bookCount,omitempty"`
	AvailableBookCount int `json:"availableBookCount,omitempty"`
	TotalBookCount     int `json:"totalBookCount,omitempty"`
	SizeOnDisk         int `json:"sizeOnDisk,omitempty"`
}
