package main

import (
	"fmt"
	"net"
)

const workerPoolSize = 20
func initTcpWorkerPool() {
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", "localhost:1234")
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
