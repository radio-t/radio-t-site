package cmd

import (
	"testing"
)

func TestUpload_Do(t *testing.T) {
	ex := &MockExecutor{}

	ex.On("Run", "mp3tags set-tags %d", 123)
	ex.On("Run", "scp %s umputun@master.radio-t.com:/srv/master-node/var/media/%s", "/tmp/rt_podcast123/rt_podcast123.mp3", "rt_podcast123.mp3")
	ex.On("Run", `ssh umputun@master.radio-t.com "chmod 644 /data/archive/radio-t/media/%s"`, "rt_podcast123.mp3")
	ex.On("Run", `ssh umputun@master.radio-t.com "find /srv/master-node/var/media -type f -mtime +60 -mtime -1200 -exec rm -vf '{}' ';'"`)
	ex.On("Run", `ssh umputun@master.radio-t.com "docker exec -i ansible /srv/deploy_radiot.sh %d"`, 123)
	ex.On("Run", `scp -P 2222 %s umputun@192.168.1.24:/data/archive.rucast.net/radio-t/media/"`, "/tmp/rt_podcast123/rt_podcast123.mp3")
	ex.On("Run", `scp %s umputun@master.radio-t.com:/data/archive/radio-t/media/%s`, "/tmp/rt_podcast123/rt_podcast123.mp3", "rt_podcast123.mp3")

	d := Upload{
		Executor: ex,
		Location: "/tmp",
	}

	d.Do(123)
	ex.AssertNumberOfCalls(t, "Run", 7)
}
