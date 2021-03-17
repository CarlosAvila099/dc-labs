package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"encoding/json"
)

type BucketData struct{
	BucketName string
	DirectoryName string
	ObjectsCount int
	DirectoriesCount int
	Extensions map[string]int
}

func printBucket(bucket BucketData){
	var message = "{\n"
	message = fmt.Sprintf(message + "	BucketName:		%s\n", bucket.BucketName)
	if bucket.DirectoryName != ""{
		message = fmt.Sprintf(message + "	DirectoryName:		%s\n", bucket.DirectoryName)
	}
	message = fmt.Sprintf(message + "	ObjectsCount:		%d\n", bucket.ObjectsCount)
	message = fmt.Sprintf(message + "	DirectoriesCount:	%d\n", bucket.DirectoriesCount)
	message = fmt.Sprintf(message + "	Extensions:{\n")
	for extension,number := range bucket.Extensions{
	message = fmt.Sprintf(message + "		%s:		%d\n", extension, number)
	}
	message = fmt.Sprintf(message + "	}\n")
	message = fmt.Sprintf(message + "}\n")
	fmt.Print(message)
}

func errorHandler(e error, i int){
	switch i{
	case 1:
		fmt.Println("There was no bucket given\n")
	case 2:
		fmt.Println("Error creating connection: ", e)
	case 3:
		fmt.Println("Error connecting: ", e)
	case 4:
		fmt.Println("Error writing information: ", e)
	case 5:
		fmt.Println("Error reading information: ", e)
	case 6:
		fmt.Println("Cant access the bucket\n")
	case 7:
		fmt.Println("Error encoding information: ", e)
	case 8:
		fmt.Println("Error decoding information: ", e)
	}
}

func handleConn(c net.Conn, m string) {
	_,err := io.WriteString(c, m)
	if err != nil{
		errorHandler(err, 4)
	}
}

func main() {
	var proxy, bucket, directory string
	flag.StringVar(&proxy, "proxy", "localhost:8000", "Connection to make")
	flag.StringVar(&bucket, "bucket", "", "bucket to enter")
	flag.StringVar(&directory, "directory", "", "directory to enter")
	flag.Parse()
	if bucket == ""{
		errorHandler(nil, 1)
		return
	}
	message := bucket
	if directory != ""{
		message += "/" + directory
	}
	conn, err := net.Dial("tcp", proxy)
	if err != nil {
		errorHandler(err, 3)
	}
	handleConn(conn, message)
	var data BucketData
	decoder := json.NewDecoder(conn)
	err = decoder.Decode(&data)
	if err != nil{
		errorHandler(err, 8)
	}
	if data.BucketName == "-1"{
		errorHandler(err, 6)
		return
	}
	printBucket(data)
	fmt.Print("\n")
	conn.Close()
}