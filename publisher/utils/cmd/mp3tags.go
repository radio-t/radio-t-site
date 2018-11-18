package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	id3 "github.com/jcs/id3-go"
	id3v2 "github.com/jcs/id3-go/v2"
	"github.com/pkg/errors"
)

// Mp3TagsCommand with command line flags and env
type Mp3TagsCommand struct {
	ShowNumber   int      `short:"s" long:"show" description:"show number"`
	File         string   `short:"f" long:"file" description:"file name"`
	Chapters     []string `short:"c" long:"chap" description:"sec::label"`
	BitRate      int      `short:"b" long:"bitrate" default:"80" description:"bitrate"`
	HugoLocation string   `long:"hugo" default:"../../hugo" description:"location of hugo"`
}

type chapter struct {
	element   string
	startSecs uint32
	endSecs   uint32
	title     string
}

var (
	reShowMdFile = regexp.MustCompile(`.*rt_podcast([0-9]*)\.mp3`)
	reChap       = regexp.MustCompile(`^- \[(.*)\]\(.*\) - \*(.*)\*`)
)

// Execute is the entry point for "mp3tags" command, called by flag parser
func (c *Mp3TagsCommand) Execute(args []string) error {
	log.Printf("mp3tags started, %+v", c)

	if c.File == "" && c.ShowNumber == 0 {
		log.Fatalf("file or show number should be defined")
	}

	fileName := c.File
	if fileName == "" {
		fileName = fmt.Sprintf("rt_podcast%d", c.ShowNumber)
	}

	showNum := c.ShowNumber
	if showNum == 0 {
		parts := reShowMdFile.FindStringSubmatch(fileName)
		if len(parts) != 2 {
			log.Fatalf("can't extract show number from %s", fileName)
		}
		var err error
		if showNum, err = strconv.Atoi(parts[1]); err != nil {
			log.Fatalf("can't extract show number from %s, %s", fileName, err)
		}
	}

	log.Printf("processing %s (%d)", fileName, showNum)
	chaps, err := c.addChapters(fileName, showNum)
	if err != nil {
		log.Fatalf("failed to add chapters, %v", err)
	}

	log.Printf("mp3-chaps completed, %d chapters", len(chaps))

	return nil
}

// based on https://github.com/jcs/mp3chap
func (c *Mp3TagsCommand) addChapters(fileName string, showNum int) ([]chapter, error) {
	mp3, err := id3.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "can't open %s", fileName)
	}

	chaps := []chapter{}
	tocChaps := []string{}

	chaps, tocChaps, err = c.mp3chaps(showNum)
	if err != nil {
		log.Printf("can't get chapters from md file for %d", showNum)
		chaps, tocChaps, err = c.manualMp3Chaps()
		if err != nil {
			return nil, errors.Wrap(err, "can't get chapters from command line")
		}
	}

	// each chapter ends where the next one starts
	for x := range chaps {
		if x < len(chaps)-1 {
			chaps[x].endSecs = chaps[x+1].startSecs
		}
	}

	finalEnd, err := c.mp3len(fileName, mp3, c.BitRate)
	if err != nil {
		return nil, errors.Wrap(err, "can't get finalEnd")
	}
	chaps[len(chaps)-1].endSecs = uint32(finalEnd * 1000)

	// ready to modify the file, clear out what's there
	mp3.DeleteFrames("CTOC")
	mp3.DeleteFrames("CHAP")

	// build a new TOC referencing each chapter
	ctocft := id3v2.V23FrameTypeMap["CTOC"]
	toc := id3v2.NewTOCFrame(ctocft, "toc", true, true, tocChaps)
	mp3.AddFrames(toc)

	// add each chapter
	chapft := id3v2.V23FrameTypeMap["CHAP"]
	for _, c := range chaps {
		ch := id3v2.NewChapterFrame(chapft, c.element, c.startSecs, c.endSecs, 0, 0, true, c.title, "", "")
		mp3.AddFrames(ch)
	}

	mp3.Close()
	return chaps, nil
}

func (c *Mp3TagsCommand) mp3len(file string, mp3 *id3.File, bitrate int) (int, error) {

	if tlenf := mp3.Frame("TLEN"); tlenf != nil {
		if tlenft, ok := tlenf.(*id3v2.TextFrame); ok {
			if tlen, err := strconv.Atoi(tlenft.Text()); err != nil {
				return tlen, nil
			}
		}
	}

	info, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	return int(info.Size()) * 8 / bitrate, nil
}

func (c *Mp3TagsCommand) mp3chaps(showNumber int) (chaps []chapter, tocChaps []string, err error) {
	mdFile := fmt.Sprintf("%s/content/posts/podcast-%d.md", c.HugoLocation, showNumber)
	mdData, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(string(mdData), "\n")
	for _, line := range lines {

		parts := reChap.FindStringSubmatch(line)
		if len(parts) != 3 {
			continue
		}
		ts, err := time.ParseInLocation("03:04:05", parts[2], time.UTC)
		if err != nil {
			return nil, nil, err
		}
		element := fmt.Sprintf("chp%d", len(chaps))
		st := ts.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))
		title := string(parts[1])
		ch := chapter{title: title, startSecs: uint32(st.Seconds() * 1000), element: element}
		log.Printf("chapter %d - %v", len(chaps), ch)
		chaps = append(chaps, ch)
		tocChaps = append(tocChaps, element)
	}
	return chaps, tocChaps, err
}

func (c *Mp3TagsCommand) manualMp3Chaps() (chaps []chapter, tocChaps []string, err error) {

	for _, c := range c.Chapters {
		elems := strings.Split(c, "::")
		if len(elems) != 2 {
			log.Fatal("incorrect chapter fromat, expecting sss::label (123:some text here)")
		}
		st, err := strconv.Atoi(elems[0])
		if err != nil {
			return nil, nil, errors.Wrapf(err, "incorrect time fromat for %s", elems[0])
		}
		element := fmt.Sprintf("chp%d", len(chaps))
		tocChaps = append(tocChaps, element)

		chap := chapter{
			element:   element,
			startSecs: uint32(st * 1000),
			endSecs:   0,
			title:     elems[1],
		}
		chaps = append(chaps, chap)
	}
	return chaps, tocChaps, nil
}
