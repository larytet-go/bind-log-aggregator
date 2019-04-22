package main

import (
	"flag"
	"fmt"
	"time"

	"transactionlogger"
)

func main() {

	useUdp := flag.Bool("udp", false, "Use UDP for rsyslog")
	loggerUrl := flag.String("logger", "", "rsyslog://127.0.0.1")
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
}
