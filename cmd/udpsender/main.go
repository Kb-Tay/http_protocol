package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const PORT = ":42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", PORT)

	if err != nil {
		log.Fatal("Failed to get UDP Addr")
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Fatal("Failed to create udp connection")
	}

	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Failed to read line")
		}

		_, err = udpConn.Write([]byte(line))

		if err != nil {
			log.Println("Faield to write" + err.Error())
		}
	}
}
