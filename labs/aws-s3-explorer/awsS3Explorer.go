package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"encoding/json"
	"golang.org/x/net/html"
	"time"
	"strings"
	"path/filepath"
)

type BucketData struct{
	BucketName string
	DirectoryName string
	ObjectsCount int
	DirectoriesCount int
	Extensions map[string]int
}

func errorHandler(e error, i int){
	switch i{
	case 1:
		fmt.Println("There was no bucket given")
	case 2:
		fmt.Println("Error creating connection: ", e)
	case 3:
		fmt.Println("Error connecting: ", e)
	case 4:
		fmt.Println("Error writing information: ", e)
	case 5:
		fmt.Println("Error reading information: ", e)
	case 6:
		fmt.Println("Cant access the bucket")
	case 7:
		fmt.Println("Error encoding information: ", e)
	case 8:
		fmt.Println("Error decoding information: ", e)
	}
}

func getData(c net.Conn) string{
	message := ""
	buffer := make([]byte, 1024)
	bytes, err := c.Read(buffer)
	if err != nil{
		errorHandler(err, 5)
	}
	if bytes > 0{
		message = string(buffer[:bytes])
	}
	return message
}

func transformData(name string) BucketData{
	var bucketName, directoryName string
	if strings.Contains(name, "/"){
		split := strings.Split(name, "/")
		bucketName = split[0]
		directoryName = split[1]
	} else{
		bucketName = name
		directoryName = ""
	}
	bucket := BucketData{bucketName, directoryName, -1, -1, make(map[string]int)}
	visited := make(map[string]int)
	url := "http://s3.amazonaws.com/" + bucketName
	data, err := http.Get(url)
	if err != nil{
		errorHandler(err, 6)
		return BucketData {"-1", "", 0, 0, make(map[string]int)}
	}
	tokenizer := html.NewTokenizer(data.Body)
	for {
		next := tokenizer.Next()
		token := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			break
		}
		switch next {
		case html.ErrorToken:
			continue
		case html.TextToken:
			dateLayout := "2017-06-30T13:36:23.000Z"
			_, err := time.Parse(dateLayout, token.Data)
			if err != nil {
				if directoryName == ""{
					extension := filepath.Ext(token.Data)
					if len(extension) > 0 {
						extension = extension[1:]
						if _, ext := bucket.Extensions[extension]; ext {
							bucket.Extensions[extension] += 1
						} else {
							bucket.Extensions[extension] = 1
						}
					}
					directory := filepath.Dir(token.Data)
					if len(directory) > 1 && directory != ".." {
						bucket.ObjectsCount += 1
						if _,in := visited[directory]; !in{
							visited[directory] = 1
							bucket.DirectoriesCount += 1
						}
					}
				}else{
					directory := filepath.Dir(token.Data)
					if len(directory) > 1 && directory != ".." && strings.Contains(directory, directoryName){
						bucket.ObjectsCount += 1
						if _,in := visited[directory]; !in{
							visited[directory] = 1
							bucket.DirectoriesCount += 1
						}
					}else{
						continue
					}
					extension := filepath.Ext(token.Data)
					if len(extension) > 0 {
						extension = extension[1:]
						if _, ext := bucket.Extensions[extension]; ext {
							bucket.Extensions[extension] += 1
						} else {
							bucket.Extensions[extension] = 1
						}
					}
				}
			}
		}
	}

	return bucket
}

func main() {
	var port string	
	flag.StringVar(&port, "port", "", "port used")
	flag.Parse()
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		errorHandler(err, 2)
	}
	for{
		conn, err := listener.Accept()
		if err != nil {
			errorHandler(err, 3)
		}
		info := getData(conn)
		encoder := json.NewEncoder(conn)
		data := transformData(info)
		err2 := encoder.Encode(data)
		if err2 != nil{
			errorHandler(err, 7)
		}
	}
}