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

// ShellExecutor is a simple wrapper to execute command within shell
type ShellExecutor struct {
	Dry bool
}

// Run makes the final command in printf style and panic on error
func (c *ShellExecutor) Run(cmd string, params ...string) {
	command := fmt.Sprintf("%s %s", cmd, strings.Join(params, " "))
	if err := c.do(command); err != nil {
		log.Printf("[ERROR] %v", err)
		os.Exit(1) // failed cmd stops everything
	}
}

// Do executes command and returns error if failed
func (c *ShellExecutor) do(cmd string) error {
	log.Printf("[INFO] execute: %s", cmd)
	if c.Dry {
		return nil
	}
	ex := exec.Command("sh", "-c", cmd)
	ex.Stdout = lgr.ToWriter(lgr.Default(), "INFO")
	ex.Stderr = lgr.ToWriter(lgr.Default(), "WARN")
	return errors.Wrapf(ex.Run(), "failed to run command: %s", cmd)
}
