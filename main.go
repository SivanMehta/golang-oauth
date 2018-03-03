package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "html/template"
  "os"
)

// constants
var (
  defaultLayout = "templates/layout.html"
	templateDir = "templates/"

  tmpls = map[string]*template.Template{}
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Serve home route")
	tmpls["home.html"].ExecuteTemplate(w, "base", map[string]interface{}{})
}

func main () {
  tmpls["home.html"] = template.Must(template.ParseFiles(templateDir + "home.html", defaultLayout))

  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler)
  http.Handle("/", r)
  port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

  log.Fatalln(http.ListenAndServe(":" + port, nil))
}
