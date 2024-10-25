[![CI/CD](https://github.com/rubenhoenle/pixelknecht/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/rubenhoenle/pixelknecht/actions/workflows/build.yml)

# Commanderer & Pixelknecht

Project Status: `IN PROGRESS - EARLY DEVELOPMENT`

This is a command-and-control server and client implementation for the "Pixelflut" game. This project consists of two components:

- The `Commanderer`, which is the command-and-control server. This server will be used to control which picture is drawn to which coordinates and to which pixelflut server.
- The `Pixelknecht`, which is the client implementation which will draw the picture recieved from the `Commanderer` to the canvas using the Pixelflut protocol. `Knecht` is a german word for `servant`.

## Motivation

Mainly getting startet with Golang, especially with asynchronous programming using channels and goroutines.

**The aim of this project is education, NOT to really write the best pixelflut client.**

## Development

This project is designed for development using the [Nix package manager](https://nix.dev/manual/nix/2.24/).

```bash
# enter dev shell
nix develop

# apply code format
nix fmt

# run the pixelknecht (pixelflut-client)
nix run

# run the commanderer (CnC-Server)
nix run .#commanderer

# start local pixelflut server for testing
nix develop
pixelnuke
```

For interacting with the _Commanderer_, there is currently no frontend available.
But there is a [Bruno](https://github.com/usebruno/bruno) collection available for interacting with the REST API.

I'm using the [pixelnuke](https://github.com/defnull/pixelflut) server for now.
