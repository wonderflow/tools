package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

	// connect to this socket
	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		// listen for reply
		fmt.Println("start to write file")
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("xxx", err)
				break
			}
			_, err = fmt.Println("new message", message)
			if err != nil {
				fmt.Println("write file error", err)
				break
			}
		}
	}()
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// send to socket
		_, err := fmt.Fprintf(conn, text+"\n")
		if err != nil {
			fmt.Println("err:", err)
			break
		}
	}
}
