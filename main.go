package main

import (
  "log"
  "fmt"
  "os"

  "net/http"
  "encoding/json"
  "io/ioutil"
  "html/template"

  "golang.org/x/oauth2"
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
)

type Config struct {
	ClientSecret string `json:"clientSecret"`
	ClientID     string `json:"clientID"`
	Secret       string `json:"secret"`
}

// constants
var (
  cfg *Config
  store *sessions.CookieStore
  oauthCfg *oauth2.Config
  tmpls = map[string]*template.Template{}
)

const (
  defaultLayout = "templates/layout.html"
  templateDir = "templates/"

  defaultConfigFile = "config/config.json"
  defaultServerCrt = "config/server.crt"
  defaultServerKey = "config/server.key"

  sessionName = "sess"
)

func loadConfig(file string) (*Config, error) {
  var config Config


  b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main () {
  tmpls["home.html"] = template.Must(template.ParseFiles(templateDir + "home.html", defaultLayout))

  store = sessions.NewCookieStore([]byte("sess"))

  var err error
	cfg, err = loadConfig(defaultConfigFile)
	if err != nil {
    fmt.Println(err)
		panic(err)
	}

  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)
  r.HandleFunc("/start", StartHandler)
  r.HandleFunc("/auth", AuthHandler)
  r.HandleFunc("/secure", SecureHandler)

  http.Handle("/", r)
  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = ":3000"
  }

  log.Println("Listening on ", port)
  err = http.ListenAndServeTLS(port, defaultServerCrt, defaultServerKey, nil)
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}
