package proc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// HugoSite represents a Hugo site.
type HugoSite struct {
	Location string
	BaseURL  string
}

// Post represents a podcast post from hugo.
type Post struct {
	CreatedAt time.Time
	URL       string
	Config    map[string]any
	Data      string
}

// ReadPosts reads all the Markdown files in the Location directory and returns a slice of Post struct containing the parsed data.
// It sorts the posts by creation date in descending order. BaseURL is used to construct the URL of the post.
func (h *HugoSite) ReadPosts() ([]Post, error) {
	log.Printf("[DEBUG] read hugo posts from %s", h.Location)
	var posts []Post
	err := filepath.Walk(h.Location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			post, err := h.parsePodcastMdFile(path)
			if err != nil {
				return err
			}
			posts = append(posts, post)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	log.Printf("[INFO] total hugo posts: %d", len(posts))
	return posts, nil
}

// parsePodcastMdFile parses a Markdown file containing podcast metadata.
// It reads the file at the given file path and extracts the podcast
// configuration in TOML format and the remaining content as data.
func (h *HugoSite) parsePodcastMdFile(filePath string) (Post, error) {
	var post Post
	var configLines []string
	file, err := os.Open(filePath)
	if err != nil {
		return post, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inConfig := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "+++" {
			inConfig = !inConfig
			continue
		}
		if inConfig {
			configLines = append(configLines, line)
		} else {
			post.Data += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return post, err
	}

	tomlData := strings.Join(configLines, "\n")
	if _, err := toml.Decode(tomlData, &post.Config); err != nil {
		return post, err
	}

	if dateStr, ok := post.Config["date"].(string); ok {
		post.CreatedAt, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err != nil {
			return post, err
		}
		post.URL = fmt.Sprintf("%s/p/%s/%s/",
			h.BaseURL, post.CreatedAt.Format("2006/01/02"), strings.TrimSuffix(filepath.Base(filePath), ".md"))
	}

	post.Data = strings.TrimPrefix(post.Data, "\n")
	post.Data = strings.TrimSuffix(post.Data, "\n")
	return post, nil
}
