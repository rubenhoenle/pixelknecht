# Commanderer & Pixelknecht

Project Status: `IN PROGRESS - EARLY DEVELOPMENT`

This is a command-and-control server and client implementation for the "Pixelflut" game. This project consists of two components:

- The `Commanderer`, which is the command-and-control server. This server will be used to control which picture is drawn to which coordinates and to which pixelflut server.
- The `Pixelknecht`, which is the client implementation which will draw the picture recieved from the `Commanderer` to the canvas using the Pixelflut protocol. `Knecht` is a german word for `servant`.

## Motivation

Mainly getting startet with Golang, especially with asynchronous programming using channels and goroutines.

The aim of this project is education, NOT to really write the best pixelflut client.
