package main

import (
	"flag"
	"fmt"
	"time"

	"transactionlogger"
)

func getOldestLogfile(logFolder string) {

}

// Pick the oldest not empty file, read the content, parse, output,
// remove the file content (keep the file itself), repeat
func tailLog(logFolder string) {

}

func main() {
	useUdp := flag.Bool("udp", false, "Use UDP for rsyslog")
	loggerUrl := flag.String("logger", "rsyslog://127.0.0.1", "For example 'stdout'")
	logFolder := flag.String("logfolder", "/var/log/named", "The bind's log folder")
	flag.Parse()

	fmt.Printf("Aggregator is staring ...\n")
	var publisher transactionlogger.Publisher
	msg := ""
	for {
		publisher, msg = transactionlogger.New(*loggerUrl, *useUdp)
		fmt.Printf("%s\n", msg)
		if publisher != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	go tailLog(*logFolder)
}
