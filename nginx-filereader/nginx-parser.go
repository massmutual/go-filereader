/*
Package nginxParser should be used as a way to parse an nginx log that is created by a kubernetes clusters ingress controller.
This package uses a regex restructer to parse the lines, so it assumes the default
kubernetes log format found here: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/ 
Example Main 
	package main

	import (
		"encoding/json"
		"fmt"
		"os"

		nginxParser "github.com/massmutual/go-filereader/nginx-filereader"
	)

	//example on how to use the nginx parser.

	func main() {
		var nignxlog nginxParser.NingxLog
		nignxlog.Unmatched = 0
		file, _ := os.Open("./test")
		position, err := nginxParser.ReadNginxLog(file, 0, &nignxlog)

		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("Got to position:", position)
		}
		out, _ := json.Marshal(nignxlog)
		fmt.Println(string(out))

	}

*/
package nginxParser

import (
	"bufio"
	"io"
	"strings"

	"github.com/alexflint/go-restructure"
)

// ReadNginxLog takes in the following objects: 
//  io.ReadSeeker - object with an already opened file to read 
//  int64 - The position you would like to start reading at 
//  *NignxLog - A pointer to an nginxlog object you would like to add lines to 
//  This method will return the last position read from that file, as well as an error 
// if there was a problem. check if error is nil before using object 
// The way this should be used is by keeping track of the returned pointer, and sleeping
// in your main some how. After a bit of time, return that pointer again to read more
// lines. this way you can get read more of the file to see if any new lines appear
func ReadNginxLog(input io.ReadSeeker, start int64, nignxlog *NginxLog) (int64, error) {
	// Set the offset of the input object to new one
	// if the number is above, it wont do anything and just
	// return an offset of 0? so error if start does not
	// equal new offset
	if newoffset, err := input.Seek(start, 0); err != nil {
		return newoffset, err
	}

	r := bufio.NewReader(input)
	pos := start
	for {
		data, readErr := r.ReadBytes('\n')
		var nginxLogline NginxLogLine
		pos += int64(len(data))
		if readErr != nil && readErr != io.EOF {
			return pos, readErr
		}

		theLogLine := strings.TrimSpace(string(data))

		if len(theLogLine) == 0 {
			if readErr == io.EOF {
				break
			} else {
				continue
			}
		} else {
			nignxlog.Lines = nignxlog.Lines + 1
		}

		matched, err := restructure.Find(&nginxLogline, theLogLine)
		if err != nil {
			return pos, err
		}
		if matched == false {
			nignxlog.Unmatched = nignxlog.Unmatched + 1
		} else {
			// This is because I cant get the regex to work where you match everything between 2 charaters like
			// (?<=\[)(.*?)(?=\]) - match charaters between [ and ] --- lookaheads and lookbehinds
			// For some reason they arnt working with this package. so instead of using them
			// im just matching [.*] then trimming those charaters here.
			// Clean this up when I figure out the look aheads
			nginxLogline.TimeLocal = strings.Trim(nginxLogline.TimeLocal, "[")
			nginxLogline.TimeLocal = strings.Trim(nginxLogline.TimeLocal, "]")
			nginxLogline.Proxy_upstream_name = strings.Trim(nginxLogline.Proxy_upstream_name, "[")
			nginxLogline.Proxy_upstream_name = strings.Trim(nginxLogline.Proxy_upstream_name, "]")
			nginxLogline.Request.Protocol = strings.Trim(nginxLogline.Request.Protocol, "\"")
			nginxLogline.Request.Endpoint = strings.Trim(nginxLogline.Request.Endpoint, " ")
			nginxLogline.Http_referer = strings.Trim(nginxLogline.Http_referer, "\"")
			nginxLogline.Http_user_agent = strings.Trim(nginxLogline.Http_user_agent, "\"")
			nignxlog.NginxLogLines = append(nignxlog.NginxLogLines, nginxLogline)
		}
		if readErr == io.EOF {
			break
		}
	}
	return pos, nil
}
