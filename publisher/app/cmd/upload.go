package cmd

import (
	"embed"
	"fmt"
	"time"

	"github.com/bogem/id3v2"
	log "github.com/go-pkgz/lgr"
)

//go:embed artifacts/*
var artifactsFS embed.FS

// Upload handles podcast upload to all destination
type Upload struct {
	Executor
	Location string
}

// Do runs uploads for given episode
func (u *Upload) Do(episodeNum int) error {
	mp3file := fmt.Sprintf("%s/rt_podcast%d/rt_podcast%d.mp3", u.Location, episodeNum, episodeNum)

	err := u.setMp3Tags(episodeNum)
	if err != nil {
		log.Printf("[WARN] can't set mp3 tags for %s, %v", mp3file, err)
	}

	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to master`, "-v", mp3file)
	u.Run("spot", "-e mp3:"+mp3file, `--task="deploy to nodes"`, "-v", mp3file)
	return nil
}

// setMp3Tags sets mp3 tags for given episode. It uses artifactsFS to read cover.jpg
func (u *Upload) setMp3Tags(episodeNum int) error {
	mp3file := fmt.Sprintf("%s/rt_podcast%d/rt_podcast%d.mp3", u.Location, episodeNum, episodeNum)
	log.Printf("[INFO] set mp3 tags for %s", mp3file)

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

	return tag.Save()
}
