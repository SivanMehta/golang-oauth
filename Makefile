server: install build serve

install:
	go get github.com/gorilla/mux
	go get github.com/gorilla/sessions

build:
	go build

serve:
	./golang-oauth

dev: build serve
