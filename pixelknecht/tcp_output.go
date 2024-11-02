package main

import (
	"fmt"
	"net"
)

const workerPoolSize = 15

func initTcpWorkerPool() {
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", "129.143.235.135:1234")
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
