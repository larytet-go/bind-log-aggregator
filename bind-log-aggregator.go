package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"transactionlogger"

	"github.com/hpcloud/tail"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	useUdp := flag.Bool("udp", false, "Use UDP for rsyslog")
	loggerUrl := flag.String("logger", "rsyslog://127.0.0.1", "For example 'stdout' or 'file://./dns_activity.log'")
	logFilename := flag.String("logfile", "", "The BIND's log filename, for example '/var/log/named/queries.log'")
	maxDepth := flag.Int("buffersize", 10*1000*1000, "Maxiumum number of log entries stored in RAM if rsyslog is down")
	syslogIp := flag.String("syslogip", "", "Open syslog connection for BIND, for example '0.0.0.0:514'")
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

	if len(*logFilename) != 0 {
		fmt.Printf("Tail -F %s\n", *logFilename)
		tailFile, err := tail.TailFile(*logFilename, tail.Config{Follow: true, ReOpen: true})
		if err != nil {
			fmt.Printf("Failed to tail log file %v\n", err)
			os.Exit(1)
		}
		for line := range tailFile.Lines {
			publisher.Push(line.Text)
		}
	}

	if len(*syslogIp) != 0 {
		fmt.Printf("Run syslog server %s\n", *syslogIp)
		channel := make(syslog.LogPartsChannel)
		handler := syslog.NewChannelHandler(channel)

		server := syslog.NewServer()
		server.SetFormat(syslog.RFC5424)
		server.SetHandler(handler)
		server.ListenUDP(*syslogIp)
		server.Boot()
		for logParts := range channel {
			fmt.Printf("%v", logParts)
		}
	}
}
