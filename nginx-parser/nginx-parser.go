package nginxParser

import (
	"bufio"
	"io"
	"strings"

	"github.com/alexflint/go-restructure"
)

type NingxLog struct {
	NginxLogLines []NginxLogLine
	Unmatched     int
}

type Request struct {
	_        struct{} `"`
	Method   string   `\w+`
	_        struct{} `\s`
	Endpoint string   `.*?\s`
	Protocol string   `.*?"`
}

type NginxLogLine struct {
	_                        struct{} `^`
	RemoteAddr               string   `(?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|-`
	_                        struct{} `\s-\s\[`
	TheRealIp                string   `(?:(?:[0-9]{1,3}\.){3}[0-9]{1,3})|-`
	_                        struct{} `\]\s-\s`
	RemoteUser               string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	TimeLocal                string   `\[.*?\]`
	_                        struct{} `\s`
	Request                  *Request
	_                        struct{} `\s`
	Status                   string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Body_bytes_sent          string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Http_referer             string   `".*?"`
	_                        struct{} `\s`
	Http_user_agent          string   `".*?"`
	_                        struct{} `\s`
	Request_length           string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Request_time             string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Proxy_upstream_name      string   `\[.*?\]`
	_                        struct{} `\s`
	Upstream_addr            string   `[a-zA-Z0-9_.-:]*`
	_                        struct{} `\s`
	Upstream_response_length string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Upstream_response_time   string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Upstream_status          string   `[a-zA-Z0-9_.-]*`
	_                        struct{} `\s`
	Req_id                   string   `[a-zA-Z0-9_.-]*`
	_                        string   `$`
}

func ReadNginxLog(input io.ReadSeeker, start int64, nignxlog *NingxLog) (int64, error) {
	if _, err := input.Seek(start, 0); err != nil {
		return 0, err
	}

	r := bufio.NewReader(input)
	pos := start
	for {
		var nginxLogline NginxLogLine
		data, err := r.ReadBytes('\n')
		pos += int64(len(data))
		if err == nil || err == io.EOF {
			if len(data) > 0 && data[len(data)-1] == '\n' {
				// length -1 to remove \n
				data = data[:len(data)-1]
			}
			if len(data) > 0 && data[len(data)-1] == '\r' {
				data = data[:len(data)-1]
			}

			matched, err := restructure.Find(&nginxLogline, string(data))
			if err != nil {
				return 0, err
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
		}
		if err != nil {
			if err != io.EOF {
				return pos, err
			}
			return pos, nil
		}
	}
	return 0, nil
}
