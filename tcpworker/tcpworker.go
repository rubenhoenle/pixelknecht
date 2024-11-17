package tcpworker

import (
	"fmt"
	"net"
	"sync"
)

func TcpWorker(wg *sync.WaitGroup, queue <-chan string, pixelflutConnectionString string) {
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
}
