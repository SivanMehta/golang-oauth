package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/google/go-github/github"
  "golang.org/x/net/context"
  "golang.org/x/oauth2"

  "log"
  "fmt"
  "net/http"
  "html/template"
  "os"
  "bytes"

  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "io/ioutil"
)

type Config struct {
	ClientSecret string `json:"clientSecret"`
	ClientID     string `json:"clientID"`
	Secret       string `json:"secret"`
}

// constants
var (
  cfg *Config
  defaultLayout = "templates/layout.html"
  templateDir = "templates/"
  tmpls = map[string]*template.Template{}

  defaultConfigFile = "config/config.json"
  defaultServerCrt = "config/server.crt"
  defaultServerKey = "config/server.key"

  store *sessions.CookieStore

	oauthCfg *oauth2.Config
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Serve home route")
  session, _ := store.Get(r, "sess")
  log.Println(session.Values["repos"])
  tmpls["home.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
  b := make([]byte, 16)
  rand.Read(b)
  session, _ := store.Get(r, "sess")

  state := base64.URLEncoding.EncodeToString(b)
  log.Println("Serve start route with auth code", state)
  session.Values["state"] = state
  session.Save(r, w)


  // form OAuth url from config files
  var buffer bytes.Buffer
  buffer.WriteString("https://github.com/login/oauth/authorize?")
  buffer.WriteString("client_id=")
  buffer.WriteString(cfg.ClientID)
  buffer.WriteString("&state=")
  buffer.WriteString(state)
  url := buffer.String()

  http.Redirect(w, r, url, 302)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
  session, err := store.Get(r, "sess")
  if err != nil {
		fmt.Fprintln(w, "could not read session")
		return
	}

  if r.URL.Query().Get("state") != session.Values["state"] {
    fmt.Fprintln(w, "no state match; possible csrf OR cookies not enabled")
    return
  }

  tkn := r.URL.Query().Get("code")

  ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tkn},
	)
	tc := oauth2.NewClient(ctx, ts)

  client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, "", nil)

  session.Values["repos"] = repos
  session.Values["accessToken"] = tkn
  session.Save(r, w)

  http.Redirect(w, r, "/", 302)
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
