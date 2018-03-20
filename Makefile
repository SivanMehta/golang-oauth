server: install build run

install:
	go get github.com/gorilla/mux
	go get github.com/gorilla/sessions

build:
	go build

run:
	./golang-oauth

dev: build run
