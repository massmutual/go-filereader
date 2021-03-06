# nginxFileReader
--
    import "."

Package nginxFileReader should be used as a way to parse an nginx log that is
created by a kubernetes clusters ingress controller. This package uses a regex
restructer to parse the lines, so it assumes the default kubernetes log format
found here:
https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/
### Example Main

    package main

    import (
    	"encoding/json"
    	"fmt"
    	"os"

    	nginxFileReader "github.com/massmutual/go-filereader/nginx-filereader"
    )

    //example on how to use the nginx parser.

    func main() {
    	var nignxlog nginxFileReader.NingxLog
    	nignxlog.Unmatched = 0
    	file, _ := os.Open("./test")
    	position, err := nginxFileReader.ReadNginxLog(file, 0, &nignxlog)

    	if err != nil {
    		fmt.Println("error:", err)
    	} else {
    		fmt.Println("Got to position:", position)
    	}
    	out, _ := json.Marshal(nignxlog)
    	fmt.Println(string(out))

    }

## Usage

#### func  ReadNginxLog

```go
func ReadNginxLog(input io.ReadSeeker, start int64, nignxlog *NginxLog) (int64, error)
```
ReadNginxLog takes in the following objects:

    io.ReadSeeker - object with an already opened file to read
    int64 - The position you would like to start reading at
    *NignxLog - A pointer to an nginxlog object you would like to add lines to
    This method will return the last position read from that file, as well as an error

if there was a problem. check if error is nil before using object The way this
should be used is by keeping track of the returned pointer, and sleeping in your
main some how. After a bit of time, return that pointer again to read more
lines. this way you can get read more of the file to see if any new lines appear

#### type NginxLog

```go
type NginxLog struct {
	NginxLogLines []NginxLogLine
	Unmatched     int
	Lines         int
}
```

NginxLog is the objject that will contain an array of parsed log lines as well
as a count of how many lines were parsed and how many lines failed to be parsed.

#### type NginxLogLine

```go
type NginxLogLine struct {
	RemoteAddr string `(?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|-`

	TheRealIp string `(?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|-`

	RemoteUser string `[a-zA-Z0-9_.-]*`

	TimeLocal string `\[.*?\]`

	Request *Request

	Status string `[a-zA-Z0-9_.-]*`

	Body_bytes_sent string `[a-zA-Z0-9_.-]*`

	Http_referer string `".*?"`

	Http_user_agent string `".*?"`

	Request_length string `[a-zA-Z0-9_.-]*`

	Request_time string `[a-zA-Z0-9_.-]*`

	Proxy_upstream_name string `\[.*?\]`

	Upstream_addr string `[a-zA-Z0-9_.-:]*`

	Upstream_response_length string `[a-zA-Z0-9_.-]*`

	Upstream_response_time string `[a-zA-Z0-9_.-]*`

	Upstream_status string `[a-zA-Z0-9_.-]*`

	Req_id string `[a-zA-Z0-9_.-]*`
}
```

NginxLogLine is the parsed out fields of the nginx log More details of all can
be found below:

    https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/

All regex fields assume you are using the default logging formating created by
kubernetes

#### type Request

```go
type Request struct {
	Method string `\w+`

	Endpoint string `.*?\s`
	Protocol string `.*?"`
}
```

Request struct is the details of the request made:

    method - GET/PUT/DELETE
    Endpoint - The endpoint that was called EX: /health
    protocol - What is the protocol used EX: HTTP/1.1

#### Generated by godocdown
To regenerate

    go get github.com/robertkrimen/godocdown/godocdown
    godocdown nginx-filereader
