package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"

  "log"
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

func main () {
  tmpls["home.html"] = template.Must(template.ParseFiles(templateDir + "home.html", defaultLayout))

  store = sessions.NewCookieStore([]byte("sess"))

  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)
  r.HandleFunc("/start", StartHandler)

  http.Handle("/", r)
  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "3000"
  }

  log.Println("Listening on port 3000")
  log.Fatalln(http.ListenAndServe(":" + port, nil))
}
