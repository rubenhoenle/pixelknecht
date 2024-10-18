package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	go draw(conn)
}

func draw(conn net.Conn) {
	for {
		for i := 0; i <= 500; i++ {
			cmd := fmt.Sprintf("PX %d %d 0000FF\n", i, i)
			_, err := conn.Write([]byte(cmd))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
