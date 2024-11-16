package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/rubenhoenle/pixelknecht/model"
)

const workerPoolSize = 15

func getPixelflutServerStringFromCommanderer() string {
	response, err := http.Get(getCommandererUrl() + "/api/server")
	if err != nil {
		fmt.Print(err.Error())
		// TODO: error handling
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// parse the response
	var server model.PixelflutServer
	err = json.Unmarshal([]byte(string(responseData)), &server)
	if err != nil {
		fmt.Println("Error:", err)
		// TODO: error handling
	}
	println(server.Host)
	println(server.Port)
	str := fmt.Sprintf("%s:%d", server.Host, server.Port)
	println(str)
	return str
}

func initTcpWorkerPool() {
	pixelflutConnectionString := getPixelflutServerStringFromCommanderer()
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", pixelflutConnectionString)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()
			for {
				cmd := <-queue
				_, err := conn.Write([]byte(cmd))
				if err != nil {
					fmt.Println(err)
					wg.Done()
					return
				}
			}
			// TODO: implement graceful shutdown
		}()
	}
}
