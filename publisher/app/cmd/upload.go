package cmd

import (
	"fmt"
	"path"

	log "github.com/go-pkgz/lgr"
)

// Upload handles podcast upload to all destination
type Upload struct {
	Executor
	Location string
}

// Do runs uploads for given episode
// panics on error
func (u *Upload) Do(episodeNum int) {
	mp3file := fmt.Sprintf("%s/rt_podcast%d/rt_podcast%d.mp3", u.Location, episodeNum, episodeNum)

	log.Printf("[INFO] set mp3 tags for %s", mp3file)
	u.Run("mp3tags set-tags %d", episodeNum)

	log.Printf("[INFO] upload %s to master.radio-t.com", mp3file)
	u.Run("scp %s umputun@master.radio-t.com:/srv/master-node/var/media/%s", mp3file, path.Base(mp3file))

	log.Printf("[INFO] set permission for %s on master.radio-t.com", mp3file)
	u.Run(`ssh umputun@master.radio-t.com "chmod 644 /data/archive/radio-t/media/%s"`, path.Base(mp3file))

	log.Printf("[INFO] remove old media files")
	u.Run(`ssh umputun@master.radio-t.com "find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'"`)

	log.Printf("[INFO] run ansible tasks")
	u.Run(`ssh umputun@master.radio-t.com "docker exec -i ansible /srv/deploy_radiot.sh %d"`, episodeNum)

	log.Printf("[INFO] copy to hp-usrv archives")
	u.Run(`scp -P 2222 %s umputun@192.168.1.24:/data/archive.rucast.net/radio-t/media/"`, mp3file)

	log.Printf("[INFO] upload to archive site")
	u.Run(`scp %s umputun@master.radio-t.com:/data/archive/radio-t/media/%s`, mp3file, path.Base(mp3file))

	return
}
