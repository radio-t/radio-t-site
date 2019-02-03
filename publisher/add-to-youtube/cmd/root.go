package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/radio-t/radio-t-site/publisher/add-to-youtube/youtube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type entry struct {
	URL        string      `json:"url"`                   // url поста
	Title      string      `json:"title"`                 // заголовок поста
	Date       time.Time   `json:"date"`                  // дата-время поста в RFC3339
	Categories []string    `json:"categories"`            // список категорий, массив строк
	Image      string      `json:"image,omitempty"`       // url картинки
	FileName   string      `json:"file_name,omitempty"`   // имя файла
	Body       string      `json:"body,omitempty"`        // тело поста в HTML
	ShowNotes  string      `json:"show_notes,omitempty"`  // пост в текстовом виде
	AudioURL   string      `json:"audio_url,omitempty"`   // url аудио файла
	TimeLabels []timeLabel `json:"time_labels,omitempty"` // массив временых меток тем
}

type timeLabel struct {
	Topic    string    `json:"topic"`              // название темы
	Time     time.Time `json:"time"`               // время начала в RFC3339
	Duration int       `json:"duration,omitempty"` // длительность в секундах
}

type siteAPIError struct {
	Message string `json:"error"`
}

func (e siteAPIError) Error() string {
	return e.Message
}

const descriptionFormat = `%s

Темы %s выпуска:

%s


Лог чата: http://chat.radio-t.com/logs/radio-t-%s.html
Аудио: %s

Информация о подкасте: https://radio-t.com/info/
Лицензия: https://radio-t.com/license/`

func makeEpisodeDescription(id string, e *entry) string {
	return fmt.Sprintf(descriptionFormat, e.URL, id, e.ShowNotes, id, e.AudioURL)
}

func getEpisodeInfo(id string) (*entry, error) {
	resp, err := http.Get(fmt.Sprintf("https://radio-t.com/site-api/podcast/%s", id))
	if err != nil {
		return nil, errSiteAPIRequest(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading request body, got: %v", err)
	}
	if resp.StatusCode != 200 {
		var e siteAPIError
		if err := json.Unmarshal(b, &e); err != nil {
			return nil, errJSONUnmarshal(err)
		}
		return nil, errSiteAPIRequest(e)
	}
	var e entry
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, errJSONUnmarshal(err)
	}
	return &e, nil
}

func download(from, to string) error {
	fmt.Printf("Start download a file `%s` to `%s`\n", from, to)

	resp, err := http.Get(from)
	if err != nil {
		return fmt.Errorf("Error downloading an audio file, got: %v", err)
	}
	defer resp.Body.Close()

	media, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("Error creation a file, got: %v", err)
	}
	defer media.Close()

	if _, err := io.Copy(media, resp.Body); err != nil {
		return fmt.Errorf("Error saving a file to fs, got: %v", err)
	}

	fmt.Printf("File `%s` downloaded\n", to)
	return nil
}

func process(id string) error {
	e, err := getEpisodeInfo(id)
	if err != nil {
		return err
	}
	d := makeEpisodeDescription(id, e)

	// prepare temprorary directory
	dir, err := ioutil.TempDir("", "add-to-youtube")
	if err != nil {
		return fmt.Errorf("Error creation a temprorary directory, got: %v", err)
	}
	defer os.RemoveAll(dir)

	filename := path.Join(dir, e.FileName+".mp3")

	if err := download(e.AudioURL, filename); err != nil {
		return err
	}

	v, err := youtube.Upload(filename, e.Title, d, "22", "podcast,radio-t", "public", viper.GetString("secrets"))
	if err != nil {
		return errJSONUnmarshal(err)
	}
	fmt.Printf("A podcast episode %s uploaded.\nhttps://youtu.be/%s\n", id, v.Id)
	return nil
}

func errSiteAPIRequest(err error) error {
	return fmt.Errorf("Error making site api request, got: %v", err)
}

func errJSONUnmarshal(err error) error {
	return fmt.Errorf("Error json unmarshaling, got: %v", err)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "add-to-youtube {episodeID}",
	Short: "Upload a radio-t podcast episode to Youtube",
	Long: `Upload a radio-t podcast episode to Youtube.

This application is a tool to generate a video file from an audio file,
then uses metadatas from site api to upload it to Youtube.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(path.Join(viper.GetString("secrets"), "youtube-secret.json")); os.IsNotExist(err) {
			fmt.Println("Before use this command you need authorize an user")
			os.Exit(1)
		}

		episodeID := args[0]
		if _, err := strconv.Atoi(episodeID); err != nil {
			fmt.Printf("{episodeID} must be a number, got: %v\n", err)
			os.Exit(1)
		}
		if err := process(episodeID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.add-to-youtube.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".add-to-youtube" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".add-to-youtube")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
