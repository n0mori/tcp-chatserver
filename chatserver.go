package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var connections []*net.TCPConn
var log *os.File

func accepter(server *net.TCPListener) {
	for {
		conn, err := server.AcceptTCP()

		if err != nil {
			return
		}

		connections = append(connections, conn)

		fmt.Printf("%s connected\n", conn.LocalAddr())

		go listener(conn)
	}
}

func listener(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')

		if err == io.EOF {
			fmt.Printf("O cliente %s se desconectou\n", conn.LocalAddr())
			return
		}

		broadcaster(str)

		fmt.Print(str)
	}

}

func broadcaster(str string) {
	for _, conn := range connections {
		if conn != nil {
			writer := bufio.NewWriter(conn)
			writer.WriteString(str)
			writer.Flush()
		}
	}
	log.WriteString(str)
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":2000")
	server, err := net.ListenTCP("tcp", addr)

	now := time.Now()
	log, err = os.OpenFile(now.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	connections = make([]*net.TCPConn, 100)

	if err != nil {
		println(err)
		os.Exit(0)
	}

	accepter(server)

	defer log.Close()

}
