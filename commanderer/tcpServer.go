package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn         net.Conn
	lastActivity time.Time
}

type TcpServer struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func RunTcpServer() {
	server := &TcpServer{clients: make(map[string]*Client)}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on", port)

	go server.checkConnections()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		server.addClient(conn)
		client := server.clients[conn.RemoteAddr().String()]
		go server.handleConnection(client)
	}
}

func (s *TcpServer) addClient(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client := &Client{conn: conn, lastActivity: time.Now()}
	s.clients[conn.RemoteAddr().String()] = client
	fmt.Printf("Added client %s\n", conn.RemoteAddr().String())
}

func (s *TcpServer) removeClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	addr := client.conn.RemoteAddr().String()
	client.conn.Close()
	delete(s.clients, addr)
	fmt.Printf("Removed client %s\n", addr)
}

func (s *TcpServer) handleConnection(client *Client) {
	defer s.removeClient(client)

	reader := bufio.NewReader(client.conn)
	for {
		client.conn.SetReadDeadline(time.Now().Add(readTimeout))

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error from %s: %v\n", client.conn.RemoteAddr(), err)
			return
		}

		if message == "heartbeat\n" {
			fmt.Printf("Received heartbeat from %s\n", client.conn.RemoteAddr())
			client.lastActivity = time.Now() // Update last activity time
		}
	}
}

func (s *TcpServer) checkConnections() {
	for {
		time.Sleep(checkInterval)
		s.mu.Lock()
		for addr, client := range s.clients {
			if time.Since(client.lastActivity) > readTimeout {
				fmt.Printf("Client %s timed out\n", addr)
				s.removeClient(client)
			}
		}
		s.mu.Unlock()
	}
}
