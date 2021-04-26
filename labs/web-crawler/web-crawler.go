// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
//
// Crawl3 adds support for depth limiting.
//
package main

import (
	"fmt"
	"os"
	"flag"
	"strconv"

	"gopl.io/ch5/links"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func errorHandler(c int, e error, l int){
	switch(c){
		case 1:
			fmt.Println("Please enter a limit bigger than " + strconv.Itoa(l))
			break
		case 2:
			fmt.Println("Please enter a name for the file")
			break
		case 3:
			fmt.Println("Please enter a URL to crawl through")
			break
		case 4:
			fmt.Println("There was an error while crawling in depth " + strconv.Itoa(l) + ": " + e.Error())
			break
		case 5:
			fmt.Println("There was an error while creating the file: " + e.Error())
			break
		case 6: 
			fmt.Println("There was an error while writing in the file: " + e.Error())
			break
	}
}

func createFile(name string, links []string){
	file, err := os.Create(name)
    if err != nil {
		errorHandler(5, err, 0)
		return
    }
	defer file.Close()
	var text string
	for _, link := range links{
		text += link
		text += "\n" 
	}
	_, err = file.WriteString(text)
	if err != nil {
		errorHandler(6, err, 0)
		return
	}
}

func crawl(url string, depth int) []string {
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		errorHandler(4, err, depth)
		return make([]string, 0)
	}
	return list
}

func depthLimiter(url string, depthLimit int, depth int, visited map[string]bool) map[string]bool{
	visited[url] = true
	if depth < depthLimit{
		list := crawl(url, depth)
		for _, link := range list{
			if !visited[link]{
				visited = depthLimiter(link, depthLimit, depth+1, visited)
			}
		}
	}
	return visited
}

func getKeys(dict map[string]bool) []string{
	keys := make([]string, 0, len(dict))
	for key := range dict{
		keys = append(keys, key)
	}
	return keys
}

//!-sema

//!+
func main() {
	var limit int
	var result string
	flag.IntVar(&limit, "depth", 0, "Depth of the web-crawler")
	flag.StringVar(&result, "results", "", "Name of file")
	flag.Parse()

	if limit < 1{
		errorHandler(1, nil, limit)
		return
	}
	if result == ""{
		errorHandler(2, nil, 0)
		return
	}
	if len(flag.Args()) == 0{
		errorHandler(3, nil, 0)
		return
	}
	url := flag.Args()[0]

	seen := make(map[string]bool)
	links := getKeys(depthLimiter(url, limit, 0, seen))
	createFile(result, links)
	fmt.Println("The crawling has ended, please check " + result + " to see the results.")
}

//!-
