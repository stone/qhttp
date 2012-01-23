#qhttp - query remote http server headers

Query servers in "paralell" (using goroutines) for http headers.

Usage:

    usage: httpid [flags] url [url...]
        -H="Server": Which header(s) to show
        -f="": read urls from file
        -n=2: number of CPU cores to use
        -v=false: verbose

Example:

    stone@ppo2:~$ ./qhttp -H="Server Expires" www.reddit.com www.lwn.net
    [0] http://www.reddit.com : 200 OK : [ '; DROP TABLE servertypes; -- ]
    [1] http://www.lwn.net : 200 OK : [ Apache | -1 ]

Note 1: You need the [go][] runtime, <http://golang.org/>

Note 2: this is just a toy project in my adventures in the go language, it probably works
but not the cleanest code around ;) 

[go]:http://golang.org/  "The Go Programming language"
