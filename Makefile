server: install build serve

install:
	go get github.com/gorilla/mux
	go get github.com/gorilla/sessions
	go get github.com/google/go-github/github
	go get golang.org/x/oauth2
	go get golang.org/x/net/context

build:
	go build

serve:
	./golang-oauth

dev: build serve
