package main

import (
  "log"
  "fmt"
  "crypto/rand"
  "net/http"
  "encoding/base64"

  "golang.org/x/oauth2"
  "golang.org/x/net/context"
  "github.com/google/go-github/github"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Serve home route")
  session, _ := store.Get(r, sessionName)
  log.Println(session.Values["repos"])
  tmpls["home.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}

func StartHandler(w http.ResponseWriter, r *http.Request) {
  b := make([]byte, 16)
  rand.Read(b)
  session, _ := store.Get(r, sessionName)

  state := base64.URLEncoding.EncodeToString(b)
  log.Println("Serve start route with auth code", state)
  session.Values["state"] = state
  session.Save(r, w)


  // form OAuth url from config files
  url := "https://github.com/login/oauth/authorize?" +
    "client_id=" + cfg.ClientID +
    "&state=" + state

  http.Redirect(w, r, url, 302)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
  session, err := store.Get(r, sessionName)
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

  log.Println(repos)

  session.Values["repos"] = repos
  session.Values["accessToken"] = tkn
  session.Save(r, w)

  http.Redirect(w, r, "/", 302)
}

func SecureHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Served secure route")
  // Get a session.
  session, err := store.Get(r, sessionName)
  log.Println(session)
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      http.Redirect(w, r, "/", 401)
      return
  }

  log.Println(session.Values["repos"])
  tmpls["home.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}
