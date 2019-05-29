package nginxParser

// NginxLog is the objject that will contain an array of parsed log lines
// as well as a count of how many lines were parsed and how many lines
// failed to be parsed.
type NginxLog struct {
	NginxLogLines []NginxLogLine
	Unmatched     int
	Lines         int
}

// Request struct is the details of the request made:
//  method - GET/PUT/DELETE
//  Endpoint - The endpoint that was called EX: /health
//  protocol - What is the protocol used EX: HTTP/1.1
type Request struct {
	_        struct{} `"`
	Method   string   `\w+`
	_        struct{} `\s`
	Endpoint string   `.*?\s`
	Protocol string   `.*?"`
}

// NginxLogLine is the parsed out fields of the nginx log
// More details of all can be found below:
//  https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/
// All regex fields assume you are using the default logging formating created by kubernetes
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
