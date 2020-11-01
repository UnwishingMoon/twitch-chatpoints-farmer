package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/nicklaw5/helix"
)

// SiteURL ...
const SiteURL = "https://www.twitch.tv/"

// App ...
type App struct {
	Client    *helix.Client
	AppTokens *helix.AccessCredentials

	Configuration *Configuration
	Settings      *Settings

	Browser *Browser
}

// Settings ...
type Settings struct {
	VideoQuality map[string]string
	VideoMuted   map[string]string
}

// NewApp ...
func NewApp() *App {
	log.Println("[INFO] Creating new App..")

	appl := &App{
		Settings: &Settings{
			VideoQuality: map[string]string{"default": "160p30"},
			VideoMuted:   map[string]string{"default": "true"},
		},
	}
	appl.readConfigFile()

	log.Println("[INFO] Creating new Twitch API Client..")

	client, err := helix.NewClient(&helix.Options{
		ClientID:     appl.Configuration.ClientID,
		ClientSecret: appl.Configuration.ClientSecret,
	})
	if err != nil {
		log.Fatalln("[FATAL] Error during client initialization.")
		log.Fatalln("[FATAL][LOG] " + err.Error())
	}

	appl.Client = client
	return appl
}

// Configuration ...
type Configuration struct {
	UserName     string `json:"username"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AuthCookie   string `json:"authCookie"`
}

func (appl *App) readConfigFile() {
	log.Println("[INFO] Reading Config File..")

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("[FATAL] Can't open config file.")
	}

	json.Unmarshal(data, &appl.Configuration)
}

// SetAppToken ...
func (appl *App) SetAppToken() {
	log.Println("[INFO] Setting Up App Token..")

	resp, err := appl.Client.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Fatalln("[FATAL] Error during App Access Token request: ", err)
	}
	appl.AppTokens = &resp.Data
	appl.Client.SetAppAccessToken(appl.AppTokens.AccessToken)
}

// GetStreams ...
func (appl *App) GetStreams() bool {
	log.Println("[INFO] Retrieving Twitch streams..")

	resp, err := appl.Client.GetStreams(&helix.StreamsParams{
		UserLogins: []string{appl.Configuration.UserName},
	})
	if err != nil {
		log.Println("[ERROR] Can't retrieve streams.")
		return false
	}

	if len(resp.Data.Streams) == 1 {
		return true
	}

	return false
}

// Browser ...
type Browser struct {
	cmd *exec.Cmd
	URL string
}

// StartBrowser ...
func (appl *App) StartBrowser() bool {
	appl.Browser = &Browser{}
	appl.Browser.cmd = exec.Command("/usr/bin/chromium-browser", "--window-size=1920,600", "--no-first-run", "--no-default-browser-check", "--headless", "--disable-gpu", "--remote-debugging-port=9222", "about:blank")

	log.Println("[INFO] Starting the browser...")

	if err := appl.Browser.cmd.Start(); err != nil {
		log.Println("[ERROR] Can't start the browser: ", err)
		return false
	}
	time.Sleep(3 * time.Second)

	log.Println("[INFO] Browser started. Yay!")

	return appl.getBrowserWsURL()
}

// CloseBrowser ...
func (appl *App) CloseBrowser() {
	if err := appl.Browser.cmd.Process.Kill(); err != nil {
		log.Fatalln("[FATAL] Can't kill Browser process: ", err)
	}
}

func (appl *App) getBrowserWsURL() bool {
	resp, err := http.Get("http://localhost:9222/json")
	if err != nil {
		log.Println("[ERROR] Can't connect to browser instance: ", err)
		appl.CloseBrowser()
		return false
	}
	defer resp.Body.Close()

	var tabs []map[string]string
	json.NewDecoder(resp.Body).Decode(&tabs)

	appl.Browser.URL = tabs[0]["webSocketDebuggerUrl"]

	return true
}
