/* golang example using http.Client on multiple goroutines

License:

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Fredrik Steen <stone4x4@gmail.com>
*/

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// Command line flags.

var (
	inputFileName = flag.String("f", "", "read urls from file")
	getHeaders    = flag.String("H", "Server", "Which header(s) to show (Default Server)")
	numCores      = flag.Int("n", runtime.NumCPU(), "number of CPU cores to use")
	verbose       = flag.Bool("v", false, "verbose")
	methodGet     = flag.Bool("get", false, "Use GET instad of HEAD")
	writeCsvFile  = flag.String("w", "", "Write to csv file")
)

// struct to hold info and results from query
type result struct {
	id      int    // Simple id
	url     string // Holds URL of query
	info    string
	headers []string
	time    time.Duration
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: httpid [flags] url [url...]")
	flag.PrintDefaults()
	os.Exit(2)
}

// Do the actual checking of the url
func geturl_head(num int, url string, c chan *result) {

	t0 := time.Now()
	var response *http.Response
	var err error
	if *methodGet {
		response, err = http.Get(url)
	} else {
		response, err = http.Head(url)
	}

	t1 := time.Now()
	time := t1.Sub(t0)

	res := &result{}
	res.id = num
	res.url = url
	res.time = time

	if err != nil {
		res.headers = nil
		if *verbose {
			res.info = err.Error()
		} else {
			res.info = "err"

		}
		c <- res
		return
	}
	defer response.Body.Close()

	//Get headers to get from flag 
	headers := strings.Split(*getHeaders, " ")
	for _, h := range headers {
		// Will be empty if no Server header is recieved
		tmphead := response.Header.Get(h)
		if tmphead == "" {
			continue
		}
		res.headers = append(res.headers, tmphead)
	}

	res.info = response.Status
	c <- res
}

// readFile returns a string array from path read from start
// to eof, removing newlines and if error returns os.Error.
func readFile(path string) (lines []string, err error) {

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for true {
		line, errr := reader.ReadString('\n')
		if errr == io.EOF {
			break
		}
		// Skip empty lines
		if line == "\n" {
			continue
		}
		lines = append(lines, line[:len(line)-1])
	}
	if err == io.EOF {
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

func NewCsv(filename string) (*csv.Writer, error) {
	fd, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	csv := csv.NewWriter(fd)
	return csv, nil
}

func writeCsvLine(w *csv.Writer, res *result) {
	headers_joined := strings.Join(res.headers, ";")
	record := []string{res.url, res.info, headers_joined, res.time.String()}
	err := w.Write(record)
	if err != nil {
		fmt.Println("Problems writing csv to files")
	}
	w.Flush()
}

func main() {
	// Handle command line args
	flag.Usage = usage
	flag.Parse()

	runtime.GOMAXPROCS(*numCores)

	var (
		urls []string
		err  error
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

	c := make(chan *result)

	for i, _ := range urls {
		furl := fixurl(urls[i])
		urls[i] = furl
		go geturl_head(i, urls[i], c)
	}

	if *writeCsvFile != "" {
		csv, err := NewCsv(*writeCsvFile)
		if err != nil {
			fmt.Println(err)
		}

		totalurls := len(urls)
		for i, _ := range urls {
			fmt.Printf("query %d of %d done\r", i, totalurls)
			res := <-c
			writeCsvLine(csv, res)
		}
		fmt.Println("")

	} else {
		for i, _ := range urls {
			res := <-c
			fmt.Printf("[%d] %s : %s : %s time=%v\n", i, urls[res.id], res.info, res.headers, res.time)
		}
	}
}
