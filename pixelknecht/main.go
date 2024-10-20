package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
    "encoding/json"
)

// TODO: check if we can use the struct from the "commanderer" module here instead of duplicating it
type floodMode struct {
	Enabled bool `json:"enabled"`
}

var wg sync.WaitGroup

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	wg.Add(1)
	go commandHandler(conn, 3)

	// wait for the goroutines to finish
	wg.Wait()
}

func generateColor() string {
	buf := make([]byte, 3)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(buf)
}

func commandHandler(conn net.Conn, pollIntervalSec int) {
	// define the initial mode
	var mode floodMode
	mode.Enabled = false

	ctx, cancel := context.WithCancel(context.Background())

	for {
		previousMode := mode
		mode = getModeFromCommanderer()

		// check if the mode changed in the meantime, if so, react to it
		if previousMode.Enabled != mode.Enabled {
			if mode.Enabled {
                fmt.Print("Starting flooding...")
	            ctx, cancel = context.WithCancel(context.Background())
				wg.Add(1)
				go draw(ctx, conn)
			} else {
                fmt.Print("Stopping flooding...")
				cancel()
			}
		}

		time.Sleep(time.Duration(pollIntervalSec) * time.Second)
	}
}

func getModeFromCommanderer() floodMode {
	response, err := http.Get("http://localhost:9000/mode")
	if err != nil {
		fmt.Print(err.Error())
        // TODO: error handling
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

    // parse the response
    var mode floodMode
	err = json.Unmarshal([]byte(string(responseData)), &mode)
	if err != nil {
		fmt.Println("Error:", err)
        // TODO: error handling
	}
    return mode
}

func draw(ctx context.Context, conn net.Conn) {
	for {
		color := generateColor()
		select {
		case <-ctx.Done(): // if cancel() execute
			wg.Done()
			return
		default:
			for i := 0; i <= 500; i++ {
				cmd := fmt.Sprintf("PX %d %d %s\n", i, i, color)
				_, err := conn.Write([]byte(cmd))
				if err != nil {
					fmt.Println(err)
					wg.Done()
					return
				}
			}
			// TODO: flood all the time and get the color from a generator
			time.Sleep(1 / 4 * time.Second)
		}
	}
}
