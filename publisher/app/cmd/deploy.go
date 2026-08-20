package cmd

import (
	"log"
)

//go:generate moq --out mocks/executor.go --pkg mocks --with-resets --skip-ensure . Executor ShellExecutorIface

// Deploy delivers site update
type Deploy struct {
	ShellExecutorIface
}

// Do commits the site to git and performs a remote site update.
func (d *Deploy) Do() {
	log.Printf("[INFO] commit site to git")
	d.RunShell(`git pull && git add . && git diff --staged --exit-code --quiet || git commit -m auto && git push`)
	log.Printf("[INFO] remote site update")
	d.Run("ssh", "umputun@master.radio-t.com", "cd /srv/radio-t/site.hugo && git pull && docker compose run --rm hugo")
}
