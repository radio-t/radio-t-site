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

// Proc handles podcast upload to all destinations. It sets mp3 tags first and then deploys to master and nodes via spot tool.
type Proc struct {
	Executor
	LocationPosts string
	Dry           bool
	SkipTransfer  bool
	Dbg           bool
}

var authors = []string{"Umputun", "Bobuk", "Gray", "Ksenks", "Alek.sys"}

// Do uploads an episode to all destinations. It takes the filename and extracts episode from this filename.
// Set all the mp3 tags and add chapters. Then deploy to master and nodes. Deploy performed by spot tool, see spot.yml
func (p *Proc) Do(mp3file string) error {
	episodeNum, err := episodeFromFile(mp3file)
	if err != nil {
		return fmt.Errorf("can't get episode number from file %s, %w", mp3file, err)
	}

	log.Printf("[INFO] process file %s, episode %d, posts location:%q", mp3file, episodeNum, p.LocationPosts)
	hugoPost := fmt.Sprintf("%s/podcast-%d.md", p.LocationPosts, episodeNum)
	log.Printf("[DEBUG] hugo post file %s", hugoPost)
	posstContent, err := os.ReadFile(hugoPost)
	if err != nil {
		return fmt.Errorf("can't read post file %s, %w", hugoPost, err)
	}
	chapters, err := p.parseChapters(string(posstContent))
	if err != nil {
		return fmt.Errorf("can't parse chapters from post %s, %w", hugoPost, err)
	}
	log.Printf("[DEBUG] chapters %v", chapters)

	err = p.setMp3Tags(mp3file, episodeNum, chapters)
	if err != nil {
		log.Printf("[WARN] can't set mp3 tags for %s, %v", mp3file, err)
	}

	if p.SkipTransfer {
		log.Printf("[WARN] skip transfer of %s", mp3file)
		return nil
	}

	newsAdminCreds := fmt.Sprintf("RT_NEWS_ADMIN:%q", os.Getenv("RT_NEWS_ADMIN"))
	args := []string{"-p /etc/spot.yml", "-e mp3:" + mp3file, "-c 2", "-v", "-e", newsAdminCreds}
	if p.Dbg {
		args = append(args, "--dbg")
	}
	p.Run("spot", args...)
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
func (p *Proc) setMp3Tags(mp3file string, episodeNum int, chapters []chapter) error {
	log.Printf("[INFO] set mp3 tags for %s", mp3file)
	if p.Dry {
		return nil
	}

	tag, err := id3v2.Open(mp3file, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("can't open mp3 file %s, %w", mp3file, err)
	}
	defer tag.Close()

	tag.DeleteAllFrames()

	tag.SetVersion(4)
	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	title := fmt.Sprintf("Радио-Т %d", episodeNum)
	tag.SetTitle(title)
	tag.SetArtist(strings.Join(authors, ", "))
	tag.SetAlbum("Радио-Т")
	tag.SetYear(fmt.Sprintf("%d", time.Now().Year()))
	tag.SetGenre("Podcast")

	// set artwork
	artwork, err := artifactsFS.ReadFile("artifacts/cover.png")
	if err != nil {
		return fmt.Errorf("can't read cover.png from artifacts, %w", err)
	}
	pic := id3v2.PictureFrame{
		MimeType:    "image/png",
		PictureType: id3v2.PTFrontCover,
		Description: "Front Cover",
		Picture:     artwork,
		Encoding:    id3v2.EncodingUTF8,
	}
	tag.AddAttachedPicture(pic)

	// we need to get mp3 duration to set the correct end time for the last chapter
	duration, err := p.getMP3Duration(mp3file)
	if err != nil {
		return fmt.Errorf("can't get mp3 duration, %w", err)
	}

	// create a CTOC frame manually
	ctocFrame := p.createCTOCFrame(chapters)
	tag.AddFrame(tag.CommonID("CTOC"), ctocFrame)

	// add other tags
	tag.AddFrame("TLEN", id3v2.TextFrame{Encoding: id3v2.EncodingUTF8, Text: strconv.FormatInt(duration.Milliseconds(), 10)})
	tag.AddFrame("TENC", id3v2.TextFrame{Encoding: id3v2.EncodingUTF8, Text: "Publisher"})

	tag.AddTextFrame(tag.CommonID("TRCK"), id3v2.EncodingUTF8, strconv.Itoa(episodeNum))
	tag.AddTextFrame(tag.CommonID("TCON"), id3v2.EncodingUTF8, "Podcast")
	tag.AddTextFrame(tag.CommonID("TCOP"), id3v2.EncodingUTF8, "Some rights reserved, Radio-T")
	tag.AddTextFrame(tag.CommonID("WXXX"), id3v2.EncodingUTF8, "https://radio-t.com")

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
			ElementID:   fmt.Sprintf("chp%d", i),
			StartTime:   chapter.Begin,
			EndTime:     endTime,
			StartOffset: id3v2.IgnoredOffset,
			EndOffset:   id3v2.IgnoredOffset,
			Title: &id3v2.TextFrame{
				Encoding: id3v2.EncodingUTF8,
				Text:     chapterTitle,
			},
		}
		tag.AddChapterFrame(chapFrame)
		log.Printf("[INFO] added chapter: %s [%v - %v]", chapterTitle, chapter.Begin, endTime)
	}

	if err := tag.Save(); err != nil {
		return err
	}
	p.ShowAllTags(mp3file)
	return nil
}

func (p *Proc) parseChapters(content string) ([]chapter, error) {
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
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "- ") {
			// Extracting the timestamp
			timestampRegex := regexp.MustCompile(`\*\s*(.*?)\s*\*`)
			timestampMatches := timestampRegex.FindStringSubmatch(line)
			if len(timestampMatches) < 2 {
				continue // Skip if no valid timestamp
			}
			begin, err := parseDuration(timestampMatches[1])
			if err != nil {
				return []chapter{}, fmt.Errorf("can't parse duration %s, %w", timestampMatches[1], err)
			}

			// Extracting and cleaning the title and URL
			titleRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
			titleMatches := titleRegex.FindStringSubmatch(line)

			var title, url string
			if len(titleMatches) >= 3 {
				title = strings.Replace(line, titleMatches[0], titleMatches[1], 1)
				url = titleMatches[2]
			} else {
				title = line
			}

			// Cleaning the title
			title = strings.TrimPrefix(title, "- ")
			title = timestampRegex.ReplaceAllString(title, "")
			title = strings.TrimSuffix(title, " - .")
			title = strings.TrimSpace(title)

			chapters = append(chapters, chapter{
				Title: title,
				URL:   url,
				Begin: begin,
			})
		}
	}

	if len(chapters) == 1 {
		return []chapter{}, nil // Return empty slice if no chapters found and only the introduction chapter is present
	}

	if len(chapters) > 0 && chapters[0].Title != "Вступление" {
		chapters = append([]chapter{{Title: "Вступление", Begin: 0}}, chapters...)
	}
	return chapters, nil
}

// getMP3Duration returns the duration of an MP3 file given its file path.
// It takes the file path as an input and returns the duration as a time.Duration and an error if any.
func (p *Proc) getMP3Duration(filePath string) (time.Duration, error) {
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

// createCTOCFrame creates a CTOC frame for the given list of chapters in the provided mp3 file.
// Making CTOC frames manually needed because id3v2 doesn't support it directly.
func (p *Proc) createCTOCFrame(chapters []chapter) *id3v2.UnknownFrame {
	var frameBody bytes.Buffer

	// TOC ID (encoded in ASCII, equivalent to "toc".encode("ascii") in Python)
	tocID := "toc"
	frameBody.WriteString(tocID)
	frameBody.WriteByte(0x00) // Null terminator for TOC ID

	// flags (0x03 for top-level and ordered chapters)
	frameBody.WriteByte(0x03)

	// number of child elements
	frameBody.WriteByte(byte(len(chapters)))

	// append child element IDs (chapter IDs) formatted as "chapter#i"
	for i := range chapters {
		elementID := fmt.Sprintf("chp%d", i)
		frameBody.WriteString(elementID)
		frameBody.WriteByte(0x00) // Null separator for IDs
	}
	// create and return an UnknownFrame with the constructed body
	return &id3v2.UnknownFrame{Body: frameBody.Bytes()}
}

// ShowAllTags shows all tags for a given mp3 file
func (p *Proc) ShowAllTags(fname string) {
	log.Printf("[DEBUG] show all tags for %s", fname)
	tag, err := id3v2.Open(fname, id3v2.Options{Parse: true})
	if err != nil {
		log.Printf("[WARN] can't open mp3 file %s, %v", fname, err)
		return
	}
	defer tag.Close()
	frames := tag.AllFrames()

	for name, frameSlice := range frames {
		if name == "APIC" {
			continue
		}

		for _, frame := range frameSlice {
			switch f := frame.(type) {
			case id3v2.ChapterFrame:
				log.Printf("[DEBUG] frame %s: ElementID:%s StartTime:%v EndTime:%v StartOffset:%v EndOffset:%v",
					name, f.ElementID, f.StartTime, f.EndTime, f.StartOffset, f.EndOffset)

				if f.Title != nil {
					log.Printf("[DEBUG] CHAP Title: Encoding:%+v Text:%s", f.Title.Encoding, f.Title.Text)
				}
				if f.Description != nil {
					log.Printf("[DEBUG] CHAP Description: Encoding:%+v Text:%s", f.Description.Encoding, f.Description.Text)
				}
			default:
				log.Printf("[DEBUG] frame %s: %+v", name, frame)
			}
		}
	}
}

// episodeFromFile takes full path to mp3 file and returns episode number
func episodeFromFile(mp3Location string) (int, error) {
	name := filepath.Base(mp3Location)
	re := regexp.MustCompile(`rt_podcast(\d+)\.mp3`)
	matches := re.FindStringSubmatch(name)
	if len(matches) != 2 {
		return 0, fmt.Errorf("can't find episode number in %s", mp3Location)
	}
	return strconv.Atoi(matches[1])
}
