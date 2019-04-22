package main

import (
	"flag"
	"fmt"

	"transactionlogger"
)

func main() {

	useUdp := flag.Bool("udp", false, "Use UDP for rsyslog")
	loggerUrl := flag.String("logger", "", "rsyslog://127.0.0.1")
	flag.Parse()

	fmt.Printf("Aggregator is staring ...\n")
	publisher, msg := transactionlogger.New(*loggerUrl, *useUdp)
	fmt.Printf("%s\n", msg)
	publisher.Push("Hello")
}
