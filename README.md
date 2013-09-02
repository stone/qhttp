#qhttp - query remote http server headers

[![Build Status](https://drone.io/github.com/stone/qhttp/status.png)](https://drone.io/github.com/stone/qhttp/latest)

Query servers in "paralell" (using goroutines) for http headers.

Usage:

    usage: qhttp [flags] url [url...]
      -H="Server": Which header(s) to show (Default Server)
      -f="": read urls from file
      -get=false: Use GET instad of HEAD
      -n=4: number of CPU cores to use
      -v=false: verbose
      -w="": Write to csv file

Example:
    
     $ ./qhttp -H="Server Expires" www.reddit.com www.lwn.net
     [0] http://www.reddit.com : 200 OK : ['; DROP TABLE servertypes; --] time=507.421ms
     [1] http://www.lwn.net : 200 OK : [Apache -1] time=1.101533s

Note 1: You need the [go][] runtime, <http://golang.org/> (weekly)

Note 2: this is just a toy project in my adventures in the go language, it probably works
but not the cleanest code around ;) 

[go]:http://golang.org/  "The Go Programming language"
