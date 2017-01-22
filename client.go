package main

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"crypto/tls"
	"fmt"
	"net"
)

var (
	config tls.Config
)

func init() {
	CA_Pool := x509.NewCertPool()
	serverCert, err := ioutil.ReadFile("./cert.pem")
	if err != nil {
		log.Fatal("Could not load server certificate!")
	}
	CA_Pool.AppendCertsFromPEM(serverCert)

	config = tls.Config{
		RootCAs: CA_Pool,
		ServerName:"127.0.0.1",
	}
}

func main() {
	var buffer = make([]byte, 1024)

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	for {
		str := readInput()
		_, err := conn.Write([]byte(str))
		if err != nil {
			return
		}
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			return
		}
		response := string(buffer[0:bytesRead])
		fmt.Print(response)
		if response == "123\n" {
			log.Println("Encrypting connection")
			doEncrypted(conn)
		}
	}
}

func doEncrypted(unenc_conn net.Conn) {
	var buffer = make([]byte, 1024)
	conn := tls.Client(unenc_conn, &config)
	err := conn.Handshake()
	if err != nil {
		log.Fatalf("tls: handshake: %s", err)
	}
	for {
		str := readInput()
		_, err := conn.Write([]byte(str))
		if err != nil {
			conn.Close()
			return
		}
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			return
		}
		response := string(buffer[0:bytesRead])
		fmt.Print(response)
		if response == "321\n" {
			log.Println("Decrypting connection")
			return
		}
	}
}

func readInput() string {
	var s string
	fmt.Scan(&s)
	return s
}
