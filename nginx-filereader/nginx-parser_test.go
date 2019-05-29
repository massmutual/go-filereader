// Package nginxFileReader should be used as a way to parse an nginx log that is created by a kubernetes clusters
// ingress controller. This package uses a regex restructer to parse the lines, so it assumes the default
// kubernetes log format found here: // https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/
package nginxFileReader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadNginxLog(t *testing.T) {
	expected1 := []byte(`{"NginxLogLines":[{"RemoteAddr":"191.168.0.1","TheRealIp":"192.168.0.1","RemoteUser":"-","TimeLocal":"23/May/2019:19:01:08 +0000","Request":{"Method":"GET","Endpoint":"/swagger/favicon-32x32.png","Protocol":"HTTP/1.1"},"Status":"200","Body_bytes_sent":"1141","Http_referer":"https://app.example.com/swagger/index.html","Http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36","Request_length":"2383","Request_time":"0.028","Proxy_upstream_name":"namespace-pythonapp-80","Upstream_addr":"192.168.0.1:5000","Upstream_response_length":"1141","Upstream_response_time":"0.032","Upstream_status":"200","Req_id":"fbdbdcdbb7a9695d8873be471dc10a58"},{"RemoteAddr":"192.168.0.1","TheRealIp":"192.168.0.1","RemoteUser":"-","TimeLocal":"23/May/2019:19:01:08 +0000","Request":{"Method":"GET","Endpoint":"/swagger/favicon-32x32.png","Protocol":"HTTP/1.1"},"Status":"200","Body_bytes_sent":"1141","Http_referer":"https://app.example.com/swagger/index.html","Http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36","Request_length":"2383","Request_time":"0.028","Proxy_upstream_name":"namespace-pythonapp-80","Upstream_addr":"192.168.0.1:5000","Upstream_response_length":"1141","Upstream_response_time":"0.032","Upstream_status":"200","Req_id":"fbdbdcdbb7a9695d8873be471dc10a58"}],"Unmatched":0,"Lines":2}`)
	expected2 := []byte(`{"NginxLogLines":[{"RemoteAddr":"192.168.0.1","TheRealIp":"192.168.0.1","RemoteUser":"-","TimeLocal":"23/May/2019:19:01:08 +0000","Request":{"Method":"GET","Endpoint":"/swagger/favicon-32x32.png","Protocol":"HTTP/1.1"},"Status":"200","Body_bytes_sent":"1141","Http_referer":"https://app.example.com/swagger/index.html","Http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36","Request_length":"2383","Request_time":"0.028","Proxy_upstream_name":"namespace-pythonapp-80","Upstream_addr":"192.168.0.1:5000","Upstream_response_length":"1141","Upstream_response_time":"0.032","Upstream_status":"200","Req_id":"fbdbdcdbb7a9695d8873be471dc10a58"}],"Unmatched":0,"Lines":1}`)
	expectednil := []byte(`{"NginxLogLines":null,"Unmatched":0,"Lines":0}`)
	expected3 := []byte(`{"NginxLogLines":[{"RemoteAddr":"192.168.0.1","TheRealIp":"192.168.0.1","RemoteUser":"-","TimeLocal":"23/May/2019:19:01:08 +0000","Request":{"Method":"GET","Endpoint":"/swagger/favicon-32x32.png","Protocol":"HTTP/1.1"},"Status":"200","Body_bytes_sent":"1141","Http_referer":"https://app.example.com/swagger/index.html","Http_user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36","Request_length":"2383","Request_time":"0.028","Proxy_upstream_name":"namespace-pythonapp-80","Upstream_addr":"192.168.0.1:5000","Upstream_response_length":"1141","Upstream_response_time":"0.032","Upstream_status":"200","Req_id":"fbdbdcdbb7a9695d8873be471dc10a58"}],"Unmatched":1,"Lines":2}`)
	var expected1Log NginxLog
	var expected2Log NginxLog
	var expected3Log NginxLog
	var expectednillog NginxLog
	json.Unmarshal(expected3, &expected3Log)
	json.Unmarshal(expected2, &expected2Log)
	json.Unmarshal(expected1, &expected1Log)
	json.Unmarshal(expectednil, &expectednillog)

	file, _ := os.Open("./test")
	type args struct {
		input    io.ReadSeeker
		start    int64
		nignxlog *NginxLog
	}
	tests := []struct {
		name         string
		args         args
		want         int64
		wantErr      bool
		wantNginxlog *NginxLog
	}{
		{
			name: "Starting At 0 pointer in file",
			args: args{
				input:    file,
				start:    0,
				nignxlog: nil,
			},
			want:         763,
			wantErr:      false,
			wantNginxlog: &expected1Log,
		},
		{
			name: "Starting At next line in log file",
			args: args{
				input:    file,
				start:    381,
				nignxlog: nil,
			},
			want:         763,
			wantErr:      false,
			wantNginxlog: &expected2Log,
		},
		{
			name: "Starting at last line",
			args: args{
				input:    file,
				start:    763,
				nignxlog: nil,
			},
			want:         763,
			wantErr:      false,
			wantNginxlog: &expectednillog,
		},
		{
			name: "give one to big for the file",
			args: args{
				input:    file,
				start:    800,
				nignxlog: nil,
			},
			want:         800,
			wantErr:      false,
			wantNginxlog: &expectednillog,
		},
		{
			name: "give a position before the start of the file",
			args: args{
				input:    file,
				start:    -1,
				nignxlog: nil,
			},
			want:         0,
			wantErr:      true,
			wantNginxlog: &expectednillog,
		},
		{
			name: "give a position in the middle of a line",
			args: args{
				input:    file,
				start:    22,
				nignxlog: nil,
			},
			want:         763,
			wantErr:      false,
			wantNginxlog: &expected3Log,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nignxlog NginxLog
			if tt.args.nignxlog == nil {
				tt.args.nignxlog = &nignxlog
			}
			got, err := ReadNginxLog(tt.args.input, tt.args.start, tt.args.nignxlog)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadNginxLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadNginxLog() = %v, want %v", got, tt.want)
				return
			}
			if cmp.Equal(&tt.wantNginxlog, &tt.args.nignxlog) != true {
				errorMessage := fmt.Sprintf("ReadNginxLog() = Nginxlog object does not match expected\nExpected:")
				expectedOutput, expectedErr := json.Marshal(&tt.wantNginxlog)
				if expectedErr != nil {
					errorMessage = fmt.Sprintf("%s %s", errorMessage, "Unable to unmarshal expected output")
				} else {
					errorMessage = fmt.Sprintf("%s %s", errorMessage, string(expectedOutput))
				}

				errorMessage = fmt.Sprintf("%s %s", errorMessage, "\nGot:")

				gotOutput, gotErr := json.Marshal(&tt.args.nignxlog)
				if gotErr != nil {
					errorMessage = fmt.Sprintf("%s %s", errorMessage, "Unable to unmarshal the object given")
				} else {
					errorMessage = fmt.Sprintf("%s %s", errorMessage, string(gotOutput))
				}
				t.Errorf(errorMessage)
			}

		})
	}
}
