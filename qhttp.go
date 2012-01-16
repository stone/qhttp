// golang example using http.Client on multiple goroutines
package main

import (
	"http"
	"fmt"
	"flag"
	"runtime"
	"os"
	"bufio"
)

// Command line flags.

var (
	inputFileName = flag.String("f", "", "read urls from file")
	getHeaders	  = flag.String("H", "Server", "Which header(s) to show")
	numCores      = flag.Int("n", 2, "number of CPU cores to use")
	verbose       = flag.Bool("v", false, "verbose")
)

// struct to hold info and results from query
type result struct {
	id     int
	url    string
	info   string
	server string
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: httpid [flags] url [url...]")
	flag.PrintDefaults()
	os.Exit(2)
}

// Do the actual checking of the url
func geturl_head(num int, url string, c chan *result) {
	response, err := http.Head(url)

	if err != nil {
		if *verbose {
			c <- &result{num, "", err.String(), ""}
		} else {
			c <- &result{num, "", "err", "err"}
		}
		return
	}

	defer response.Body.Close()

	// Will be empty if no Server header is recieved
	server := response.Header.Get("Server")
	c <- &result{num, url, response.Status, server}
}

// readFile returns a string array from path read from start
// to eof, removing newlines and if error returns os.Error.
func readFile(path string) (lines []string, err os.Error) {

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for true {
		line, errr := reader.ReadString('\n')
		if errr == os.EOF {
			break
		}
		// Skip empty lines
		if line == "\n" {
			continue
		}
		lines = append(lines, line[:len(line)-1])
	}
	if err == os.EOF {
		err = nil
	}
	return
}

// fixurl pretty stupid function that checks if
// url starts with http, if not appends http://
func fixurl(url string) (furl string) {
	if url[:4] != "http" {
		furl = "http://" + url
	} else {
		furl = url
	}
	return
}

func main() {
	// Handle command line args
	flag.Usage = usage
	flag.Parse()

	runtime.GOMAXPROCS(*numCores)

	var (
		urls []string
		err  os.Error
	)

	// if we got a file to read urls from
	if *inputFileName != "" {
		urls, err = readFile(*inputFileName)
		if err != nil {
			fmt.Printf("\nOpen Error => %s\n\n", err)
			os.Exit(1)
		}
	} else {
		// if we got args
		if flag.NArg() > 0 {
			urls = flag.Args()
		} else {
			usage()
		}
	}

	c := make(chan *result, 100)

	for i, _ := range urls {
		furl := fixurl(urls[i])
		urls[i] = furl
		go geturl_head(i, urls[i], c)
	}

	for i, _ := range urls {
		res := <-c
		fmt.Printf("[%d] %s : %s : %s\n", i, urls[res.id], res.info, res.server)
	}
}