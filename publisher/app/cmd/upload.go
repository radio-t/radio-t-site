package cmd

import (
	"embed"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bogem/id3v2/v2"
	log "github.com/go-pkgz/lgr"
)

//go:embed artifacts/*
var artifactsFS embed.FS

// Upload handles podcast upload to all destinations. It uses spot to deploy and set mp3 tags before deploy
type Upload struct {
	Executor
	LocationMp3   string
	LocationPosts string
	Dry           bool
}

// Do runs uploads for given episode
func (u *Upload) Do(episodeNum int) error {
	log.Printf("[INFO] upload episode %d, mp3 location:%q, posts location:%q", episodeNum, u.LocationMp3, u.LocationPosts)
	mp3file := fmt.Sprintf("%s/rt_podcast%d/rt_podcast%d.mp3", u.LocationMp3, episodeNum, episodeNum)
	log.Printf("[DEBUG] mp3 file %s", mp3file)
	hugoPost := fmt.Sprintf("%s/podcast-%d.md", u.LocationPosts, episodeNum)
	log.Printf("[DEBUG] hugo post file %s", hugoPost)
	posstContent, err := os.ReadFile(hugoPost)
	if err != nil {
		return fmt.Errorf("can't read post file %s, %w", hugoPost, err)
	}
	chapters, err := u.parseChapters(string(posstContent))
	if err != nil {
		return fmt.Errorf("can't parse chapters from post %s, %w", hugoPost, err)
	}
	log.Printf("[DEBUG] chapters %v", chapters)

	err = u.setMp3Tags(episodeNum, chapters)
	if err != nil {
		log.Printf("[WARN] can't set mp3 tags for %s, %v", mp3file, err)
	}

	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to master`, "-v", mp3file)
	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to nodes"`, "-v", mp3file)
	return nil
}

// chapter represents a single chapter in the podcast
type chapter struct {
	Title string
	URL   string
	Begin time.Duration
}

// setMp3Tags sets mp3 tags for given episode. It uses artifactsFS to read cover.jpg
func (u *Upload) setMp3Tags(episodeNum int, chapters []chapter) error {
	mp3file := fmt.Sprintf("%s/rt_podcast%d/rt_podcast%d.mp3", u.LocationMp3, episodeNum, episodeNum)
	log.Printf("[INFO] set mp3 tags for %s", mp3file)
	if u.Dry {
		return nil
	}
	tag, err := id3v2.Open(mp3file, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("can't open mp3 file %s, %w", mp3file, err)
	}
	defer tag.Close()

	tag.SetTitle(fmt.Sprintf("Радио-Т %d", episodeNum))
	tag.SetArtist("Umputun, Bobuk, Gray, Ksenks, Alek.sys")
	tag.SetAlbum("Радио-Т")
	tag.SetYear(fmt.Sprintf("%d", time.Now().Year()))
	tag.SetGenre("Podcast")
	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	artwork, err := artifactsFS.ReadFile("artifacts/cover.png")
	if err != nil {
		return fmt.Errorf("can't read cover.jpg from artifacts, %w", err)
	}

	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/png",
		PictureType: id3v2.PTFrontCover,
		Description: "Front Cover",
		Picture:     artwork,
	}
	tag.AddAttachedPicture(pic)

	for i, chapter := range chapters {
		var endTime time.Duration
		if i < len(chapters)-1 {
			// use the start time of the next chapter as the end time
			endTime = chapters[i+1].Begin
		} else {
			endTime = 0
		}
		chapFrame := id3v2.ChapterFrame{
			ElementID:   strconv.Itoa(i + 1),
			StartTime:   chapter.Begin,
			EndTime:     endTime,
			StartOffset: id3v2.IgnoredOffset,
			EndOffset:   id3v2.IgnoredOffset,
			Title: &id3v2.TextFrame{
				Encoding: id3v2.EncodingUTF8,
				Text:     chapter.Title,
			},
			Description: &id3v2.TextFrame{
				Encoding: id3v2.EncodingUTF8,
				Text:     chapter.Title,
			},
		}
		tag.AddChapterFrame(chapFrame)
	}

	return tag.Save()
}

// parseChapters parses the input content and returns a slice of chapters
func (u *Upload) parseChapters(content string) ([]chapter, error) {
	parseDuration := func(timestamp string) (time.Duration, error) {
		parts := strings.Split(timestamp, ":")
		if len(parts) != 3 {
			return 0, fmt.Errorf("invalid timestamp format")
		}

		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}

		return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
	}

	chapters := []chapter{}
	chapterRegex := regexp.MustCompile(`-\s+\[(.*?)\]\((.*?)\)\s+-\s+\*(.*?)\*\.`)

	matches := chapterRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 4 {
			title := match[1]
			url := match[2]
			timestamp := match[3]

			begin, err := parseDuration(timestamp)
			if err != nil {
				return nil, err
			}

			chapters = append(chapters, chapter{
				Title: title,
				URL:   url,
				Begin: begin,
			})
		}
	}

	return chapters, nil
}
