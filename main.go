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
