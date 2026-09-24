package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-pkgz/lgr"
	"github.com/pkg/errors"
)

// Executor is a simple interface to run commands
type Executor interface {
	Run(cmd string, params ...string)
}

// ShellExecutorIface runs commands and shell scripts, for callers which need both
type ShellExecutorIface interface {
	Executor
	RunShell(script string)
}

// LastShow get the number of the latest published podcast via site-api
// GET /last/{posts}?categories=podcast
func LastShow(client http.Client, siteAPI string) (int, error) {
	resp, err := client.Get(fmt.Sprintf("%s/last/1?categories=podcast", siteAPI))
	if err != nil {
		return -1, errors.Wrap(err, "can't get last shows")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return -1, errors.Errorf("invalid status code %s", resp.Status)
	}

	//noinspection GoPreferNilSlice
	showInfo := []struct {
		Num int `json:"show_num"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&showInfo); err != nil {
		return -1, errors.Wrap(err, "can't read and decode")
	}

	if len(showInfo) < 1 {
		return -1, errors.New("list of podcasts is empty")
	}

	return showInfo[0].Num, nil
}

// ShellExecutor runs commands, either as argv via Run or as a shell script via RunShell
type ShellExecutor struct {
	Dry bool
}

// Run executes cmd with params passed as separate arguments, without a shell, and exits on error.
// Values are never re-parsed, so paths and credentials may contain spaces and shell metacharacters.
func (c *ShellExecutor) Run(cmd string, params ...string) {
	log.Printf("[INFO] execute: %s %s", cmd, strings.Join(params, " ")) //nolint:gosec // the command is composed locally, log injection is not a concern for a CLI tool
	if c.Dry {
		return
	}
	ex := exec.Command(cmd, params...) //nolint:gosec // the command is composed by the caller, not by external input
	if err := c.exec(ex, cmd); err != nil {
		log.Printf("[ERROR] %v", err)
		os.Exit(1) // failed cmd stops everything
	}
}

// RunShell executes script with sh, for commands which need shell syntax, and exits on error.
func (c *ShellExecutor) RunShell(script string) {
	log.Printf("[INFO] execute: %s", script) //nolint:gosec // the script is composed locally, log injection is not a concern for a CLI tool
	if c.Dry {
		return
	}
	ex := exec.Command("sh", "-c", script) //nolint:gosec // running a shell script is what this method is for
	if err := c.exec(ex, script); err != nil {
		log.Printf("[ERROR] %v", err)
		os.Exit(1) // failed cmd stops everything
	}
}

// exec runs the prepared command, piping its output to the log
func (c *ShellExecutor) exec(ex *exec.Cmd, descr string) error {
	ex.Stdout = lgr.ToWriter(lgr.Default(), "INFO")
	ex.Stderr = lgr.ToWriter(lgr.Default(), "WARN")
	return errors.Wrapf(ex.Run(), "failed to run command: %s", descr)
}
