# golang-oauth
Seeing if I can make a basic web app that requires login with GitHub

started based off [this](https://github.com/andrewtian/golang-github-oauth-example/blob/master/main.go)
example but we'll see how far this goes

# Getting started

The is the first time I'm using a `Makefile` to do anything, so IDK if it's correct at all. You will need a `config.json` with the correct credentials. I have included an example (with fake credentials), in `config/config.example.json`

## `make install`

Just does `go get` on all the relevant dependencies

## `make build`

Compiles `.go` files into a executable

## `make serve`

Runs executable built in `make build`

## `make dev`

Used for development, rebuilds and serves the server-side code

## `make`

If you trying to start from scratch, this will install the dependencies, re-build the server, and serve it.
