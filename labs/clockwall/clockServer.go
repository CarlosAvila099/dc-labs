// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"os"
	"io"
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn, ch chan string, timeZone string) {
	defer c.Close()
	for {
		time.Sleep(time.Second)
		message := timeZone + "\t:  " + <-ch
		io.WriteString(c, message)
	}
}

func clock(chIn chan string, chOut chan string){
	name := <- chIn
	for {
		t, err := timeIn(time.Now(), name)
		if err != nil{
			return
		}
		chOut <- t.Format("15:04:05\n")
	}
}

func timeIn(t time.Time, name string) (time.Time, error) {
    loc, err := time.LoadLocation(name)
    if err == nil {
        t = t.In(loc)
    }
    return t, err
}

func main() {
	var port string
	flag.StringVar(&port, "port", "", "port to be opened")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	ch1 := make(chan string, 1)
	ch2 := make(chan string)
	TZ := os.Getenv("TZ")
	ch1 <- TZ
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go clock(ch1, ch2)
		go handleConn(conn, ch2, TZ) // handle connections concurrently
	}
	close(ch1)
	close(ch2)
}
