package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"

  "log"
  "fmt"
  "net/http"
  "html/template"
  "os"

  "crypto/rand"
  "encoding/base64"
)

// constants
var (
  defaultLayout = "templates/layout.html"
  templateDir = "templates/"

  tmpls = map[string]*template.Template{}

  store *sessions.CookieStore

	oauthCfg *oauth2.Config
)



func HomeHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Serve home route")
  tmpls["home.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}

func StartHandler(w http.ResponseWriter, r *http.Request) {

  b := make([]byte, 16)
  rand.Read(b)
  session, _ := store.Get(r, "sess")

  state := base64.URLEncoding.EncodeToString(b)
  log.Println("Serve start route with state", state)
  session.Values["state"] = state
  session.Save(r, w)

  http.Redirect(w, r, "/", 302)
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

  tkn, err := oauthCfg.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
  if err != nil {
    fmt.Fprintln(w, "there was an issue getting your token")
    return
  }

  if !tkn.Valid() {
		fmt.Fprintln(w, "retreived invalid token")
		return
	}

  client := github.NewClient(oauthCfg.Client(oauth2.NoContext, tkn))
  ctx := r.Context()

  user, _,  err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Println(w, "error getting name")
		return
	}

  session.Values["name"] = user.Name
  session.Values["accessToken"] = tkn.AccessToken
  session.Save(r, w)

  http.Redirect(w, r, "/", 302)

}

func main () {
  tmpls["home.html"] = template.Must(template.ParseFiles(templateDir + "home.html", defaultLayout))

  store = sessions.NewCookieStore([]byte("sess"))

  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)
  r.HandleFunc("/start", StartHandler)
  r.HandleFunc("/auth", AuthHandler)

  http.Handle("/", r)
  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "3000"
  }

  log.Println("Listening on port 3000")
  log.Fatalln(http.ListenAndServe(":" + port, nil))
}
