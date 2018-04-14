package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

var connections []*net.TCPConn

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
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":2000")
	server, err := net.ListenTCP("tcp", addr)

	connections = make([]*net.TCPConn, 100)

	if err != nil {
		println(err)
		os.Exit(0)
	}

	accepter(server)

}
