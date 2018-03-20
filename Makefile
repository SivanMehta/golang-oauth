server: install build serve

install:
	go get github.com/gorilla/mux
	go get github.com/gorilla/sessions
	go get golang.org/x/oauth2
	go get github.com/google/go-github/github

build:
	go build

serve:
	./golang-oauth

dev: build serve
