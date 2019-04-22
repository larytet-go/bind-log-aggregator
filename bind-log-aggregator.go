package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"transactionlogger"

	"github.com/hpcloud/tail"
)

func main() {
	useUdp := flag.Bool("udp", false, "Use UDP for rsyslog")
	loggerUrl := flag.String("logger", "rsyslog://127.0.0.1", "For example 'stdout'")
	logFilename := flag.String("logfile", "/var/log/named/queries.log", "The bind's log filename")
	maxDepth := flag.Int("buffersize", 10*1000*1000, "Maxiumum number of log entries stored in RAM if rsyslog is down")
	flag.Parse()

	fmt.Printf("Aggregator is staring ...\n")
	var publisher transactionlogger.Publisher
	msg := ""
	for {
		publisher, msg = transactionlogger.New(*maxDepth, *loggerUrl, *useUdp)
		fmt.Printf("%s\n", msg)
		if publisher != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	tailFile, err := tail.TailFile(*logFilename, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		fmt.Printf("Failed to tail log file %v\n", err)
		os.Exit(1)
	}
	for line := range tailFile.Lines {
		fmt.Println(line.Text)
	}
}
