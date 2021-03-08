package main

import (
	"fmt"
	"os"
	"strings"
	"net"
)

func getData(c net.Conn, ch chan string){
	defer c.Close()
	buffer := make([]byte, 1024)
	bytes, err := c.Read(buffer)
	if err != nil{
		fmt.Println("Error reading connection:", err.Error())
		os.Exit(2)
	}
	if bytes > 0{
		ch <- string(buffer[:])
	}
}

func main() {
	connection := make([]string, len(os.Args[1:]))
	channel := make(chan string, len(os.Args[1:]))
	for num, part := range os.Args[1:] {
		split := strings.Split(part, "=")
		connection[num] = split[1]
	}
	for _, con := range connection{
		conn, err := net.Dial("tcp", con)
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		}
		go getData(conn, channel)
	}
	for line := range channel{
		fmt.Printf("\r%v", line)
	}
	close(channel)
}