package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bogem/id3v2/v2"
	"github.com/tcolgate/mp3"
)

//go:embed artifacts/*
var artifactsFS embed.FS

// Upload handles podcast upload to all destinations. It sets mp3 tags first and then deploys to master and nodes via spot tool.
type Upload struct {
	Executor
	LocationMp3   string
	LocationPosts string
	Dry           bool
	SkipTransfer  bool
}

// Do uploads an episode to all destinations. It takes an episode number as input and returns an error if any of the actions fail.
// It performs the following actions:
//  1. Set mp3 tags.
//  2. Deploy to master.
//  3. Deploy to nodes.
//
// deploy performed by spot tool, see spot.yml
func (u *Upload) Do(episodeNum int) error {
	log.Printf("[INFO] upload episode %d, mp3 location:%q, posts location:%q", episodeNum, u.LocationMp3, u.LocationPosts)
	mp3file := filepath.Join(u.LocationMp3, fmt.Sprintf("rt_podcast%d", episodeNum), fmt.Sprintf("rt_podcast%d.mp3", episodeNum))
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

	if u.SkipTransfer {
		log.Printf("[WARN] skip transfer of %s", mp3file)
		return nil
	}

	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to master"`, "-v")
	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to nodes"`, "-v")
	return nil
}

// chapter represents a single chapter in the podcast
type chapter struct {
	Title string
	URL   string
	Begin time.Duration
}

// setMp3Tags sets mp3 tags for a given episode. It uses artifactsFS to read cover.jpg
// and uses the chapter information to set the chapter tags.
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

	tag.DeleteAllFrames() // clear all existing tags

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	title := fmt.Sprintf("Радио-Т %d", episodeNum)
	tag.SetTitle(title)
	tag.SetArtist("Umputun, Bobuk, Gray, Ksenks, Alek.sys")
	tag.SetAlbum("Радио-Т")
	tag.SetYear(fmt.Sprintf("%d", time.Now().Year()))
	tag.SetGenre("Podcast")

	// set artwork
	artwork, err := artifactsFS.ReadFile("artifacts/cover.png")
	if err != nil {
		return fmt.Errorf("can't read cover.png from artifacts, %w", err)
	}
	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/png",
		PictureType: id3v2.PTFrontCover,
		Description: "Front Cover",
		Picture:     artwork,
	}
	tag.AddAttachedPicture(pic)

	// we need to get mp3 duration to set the correct end time for the last chapter
	duration, err := u.getMP3Duration(mp3file)
	if err != nil {
		return fmt.Errorf("can't get mp3 duration, %w", err)
	}

	// create a CTOC frame manually
	ctocFrame := u.createCTOCFrame(chapters)
	tag.AddFrame(tag.CommonID("CTOC"), ctocFrame)

	// add chapters
	for i, chapter := range chapters {
		var endTime time.Duration
		if i < len(chapters)-1 {
			endTime = chapters[i+1].Begin
		} else {
			endTime = duration
		}
		chapterTitle := chapter.Title
		if !utf8.ValidString(chapterTitle) {
			return fmt.Errorf("chapter title contains invalid UTF-8 characters")
		}
		chapFrame := id3v2.ChapterFrame{
			ElementID:   strconv.Itoa(i + 1),
			StartTime:   chapter.Begin,
			EndTime:     endTime,
			StartOffset: id3v2.IgnoredOffset,
			EndOffset:   id3v2.IgnoredOffset,
			Title: &id3v2.TextFrame{
				Encoding: id3v2.EncodingUTF8,
				Text:     chapterTitle,
			},
			Description: &id3v2.TextFrame{
				Encoding: id3v2.EncodingUTF8,
				Text:     chapterTitle,
			},
		}
		tag.AddChapterFrame(chapFrame)
	}

	return tag.Save()
}

// parseChapters parses md post content and returns a list of chapters
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

	chapters := []chapter{
		{Title: "Вступление", Begin: 0},
	}

	// get form md like this "- [Chapter One](http://example.com/one) - *00:01:00*."
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
	if len(chapters) == 1 {
		return []chapter{}, nil // no chapters found, don't return the introduction chapter
	}
	return chapters, nil
}

func (u *Upload) getMP3Duration(filePath string) (time.Duration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	d := mp3.NewDecoder(file)
	var f mp3.Frame
	var skipped int
	var duration float64

	for err == nil {
		if err = d.Decode(&f, &skipped); err != nil && err != io.EOF {
			log.Printf("[WARN] can't get duration for provided stream: %v", err)
			return 0, nil
		}
		duration += f.Duration().Seconds()
	}
	return time.Second * time.Duration(duration), nil
}

func (u *Upload) createCTOCFrame(chapters []chapter) *id3v2.UnknownFrame {
	var frameBody bytes.Buffer
	frameBody.WriteByte(0x03)                // write flags (e.g., 0x03 for top-level and ordered chapters)
	frameBody.WriteByte(byte(len(chapters))) // write the number of child elements

	// append child element IDs (chapter IDs)
	for i, _ := range chapters {
		elementID := fmt.Sprintf("%d", i+1)
		frameBody.WriteString(elementID)
		frameBody.WriteByte(0x00) // Null separator for IDs
	}

	// create and return an UnknownFrame with the constructed body
	return &id3v2.UnknownFrame{Body: frameBody.Bytes()}
}
