package proc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/russross/blackfriday/v2"
)

//go:generate moq --out mocks/http_client.go --pkg mocks --with-resets --skip-ensure . HttpClient

// RSSGenerator is a primary struct for RSS generation.
type RSSGenerator struct {
	Client          HttpClient
	BaseURL         string // BaseURL is a base URL for posts https://radio-t.com
	BaseArchiveURL  string // BaseArchiveURL is a base URL for mp3 files http://archive.rucast.net/radio-t/media"
	BaseCdnURL      string // BaseCdnURL is a base URL for mp3 files https://cdn.rucast.net"
	RssRootLocation string // RssRootLocation is a root location for rss files /srv/hugo/public
}

// FeedConfig represents a feed configuration for generating an RSS feed.
type FeedConfig struct {
	Name            string
	Title           string
	Image           string
	Count           int
	Size            bool
	FeedSubtitle    string
	FeedDescription string
}

// ItemData is the struct for each item in the feed
type ItemData struct {
	Title        string
	Description  string
	URL          string
	GUID         string
	Date         string
	Summary      string
	Image        string
	EnclosureURL string
	FileSize     int
}

// FeedData is the struct that matches the placeholders in the rssTemplate
type FeedData struct {
	FeedTitle       string
	FeedURL         string
	FeedSubtitle    string
	FeedDescription string
	FeedImage       string
	Items           []ItemData
}

// HttpClient is an interface that represents an HTTP client, compatible with the standard library.
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MakeFeed generates an RSS feed based on the given feed configuration and posts.
// It prepares the data for the RSS template by iterating through the posts and populating the feedData structure.
// Then, it parses and executes the RSS template using the feedData.
func (g *RSSGenerator) MakeFeed(feedCfg FeedConfig, posts []Post) (FeedData, error) {
	log.Printf("[DEBUG] make feed %q", feedCfg.Name)
	// preparing data for the template
	feedData := FeedData{
		FeedTitle:       feedCfg.Title,
		FeedURL:         g.BaseURL,
		FeedSubtitle:    feedCfg.FeedSubtitle,
		FeedDescription: feedCfg.FeedDescription,
		FeedImage:       feedCfg.Image,
		Items:           make([]ItemData, 0),
	}

	for _, post := range posts {
		if len(feedData.Items) >= feedCfg.Count {
			break
		}
		if _, ok := post.Config["categories"]; !ok {
			continue
		}

		categories, ok := post.Config["categories"].([]any)
		if !ok || !contains(categories, "podcast") {
			continue
		}

		// populate ItemData for each post
		item, err := g.createItemData(feedCfg, post)
		if err != nil {
			return FeedData{}, fmt.Errorf("error creating item data: %v", err)
		}
		feedData.Items = append(feedData.Items, item)
	}
	log.Printf("[INFO] total items in feed %q: %d", feedCfg.Name, len(feedData.Items))
	return feedData, nil
}

// Save parses the RSS template and execute it using the given feedData.
// The resulting RSS feed is saved to the given file path.
func (g *RSSGenerator) Save(feedCfg FeedConfig, data FeedData) error {
	savePath := filepath.Join(g.RssRootLocation, feedCfg.Name+".rss")
	log.Printf("[INFO] save feed %s to %s", feedCfg.Name, savePath)

	tmpl, err := template.New("rss").Parse(rssTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return nil

}

// createItemData prepares the ItemData for each post.
func (g *RSSGenerator) createItemData(feed FeedConfig, post Post) (ItemData, error) {
	date := post.CreatedAt.Format("Mon, 02 Jan 2006 15:04:05 MST")
	filesize := 0
	if feed.Size { // feed.Size indicates mp3 file size should be included in the feed
		var err error
		filesize, err = g.getMp3Size(post.Config["filename"].(string) + ".mp3")
		if err != nil {
			return ItemData{}, fmt.Errorf("error getting mp3 size: %v", err)
		}
	}

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{Flags: blackfriday.UseXHTML})
	data := blackfriday.Run([]byte(post.Data), blackfriday.WithRenderer(renderer))
	postDescriptionHTML := string(data)
	postDescriptionHTML = strings.TrimSuffix(postDescriptionHTML, "\n")
	rssDescriptionHTML := strings.Replace(postDescriptionHTML, "<ul>", "<p><em>Темы</em><ul>", 1)
	rssDescriptionHTML = strings.Replace(rssDescriptionHTML, "</ul>", "</ul></p>", 1)
	rssDescriptionHTML = strings.TrimSuffix(rssDescriptionHTML, "\n")

	fixedURL := post.URL
	fixedURL = strings.Replace(fixedURL, "//p", "/p", 1)
	guid := strings.Replace(fixedURL, "/podcast-", "//podcast-", 1) // to match the old feed guid
	res := ItemData{
		Description:  rssDescriptionHTML,
		URL:          fixedURL,
		GUID:         guid,
		Date:         date,
		Summary:      postDescriptionHTML,
		EnclosureURL: fmt.Sprintf("%s/%s.mp3", g.BaseCdnURL, post.Config["filename"]),
		FileSize:     filesize,
	}

	if r, ok := post.Config["title"].(string); ok {
		res.Title = r
	}
	if r, ok := post.Config["image"].(string); ok {
		res.Image = r
	}

	return res, nil
}

// getMp3Size returns the size of remote mp3 file in bytes
func (g *RSSGenerator) getMp3Size(mp3File string) (int, error) {
	url := g.BaseArchiveURL + "/" + mp3File
	if strings.HasSuffix(url, "/") {
		url = g.BaseArchiveURL + mp3File
	}
	req, err := http.NewRequest("HEAD", url, http.NoBody)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error getting response for %s: %v", req.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// check if Content-Length header is available
		if size, ok := resp.Header["Content-Length"]; ok && len(size) > 0 {
			sizeBytes, err := strconv.Atoi(size[0])
			if err != nil {
				return 0, fmt.Errorf("error converting content length %s to int: %v", size[0], err)
			}
			return sizeBytes, nil
		}
	}

	// if no Content-Length header from HEAD request, try GET request
	req, err = http.NewRequest("GET", url, http.NoBody)
	resp, err = g.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error getting response for %s: %v", req.URL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, req.URL)
	}
	return int(resp.ContentLength), nil
}

func contains(slice []any, item any) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
