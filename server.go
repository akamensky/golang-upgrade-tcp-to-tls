package main

import (
	"log"
	"crypto/tls"
	"net"
)

var (
	config tls.Config
)

func init() {
	cert, _ := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	config = tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

func main() {
	listener, _ := net.Listen("tcp", "127.0.0.1:8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("server: accepted from %s", conn.RemoteAddr())
	var buffer = make([]byte, 1024)
	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			log.Printf("server: closing from %s", conn.RemoteAddr())
			break
		}
		response := string(buffer[0:bytesRead])
		conn.Write([]byte(response + "\n"))
		if response == "123" {
			log.Printf("server: encrypting connection from %s", conn.RemoteAddr())
			handleTLSConnection(conn)
		}
	}
	conn.Close()
}

func handleTLSConnection(unenc_conn net.Conn) {
	conn := tls.Server(unenc_conn, &config)
	var buffer = make([]byte, 1024)
	conn.Handshake()
	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		response := string(buffer[0:bytesRead])
		conn.Write([]byte(response + "\n"))
		if response == "321" {
			log.Printf("server: UNencrypting connection from %s", conn.RemoteAddr())
			return
		}
	}
}
