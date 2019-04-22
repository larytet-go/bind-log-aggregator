package main

import (
	"bufio"
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
	logFilename := flag.String("logfile", "", "The BIND's log filename, for example '/var/log/named/queries.log', 'stdin'")
	maxDepth := flag.Int("buffersize", 10*1000*1000, "Maxiumum number of log entries stored in RAM if rsyslog is down")
	syslogIp := flag.String("syslogip", "", "Open syslog connection for BIND, for example '0.0.0.0:514'")
	flag.Parse()

	fmt.Printf("Aggregator is staring ...\n")
	var publisher transactionlogger.Publisher
	msg := ""
	for {
		publisher, msg = transactionlogger.Ne w(*maxDepth, *loggerUrl, *useUdp)
		fmt.Printf("%s\n", msg)
		if publisher != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if *logFilename  == "stdin" {
		fmt.Printf("Waiting for stdin\n")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			publisher.Push(line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Fail to read stdin %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
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
		os.Exit(0)
	}

	if len(*syslogIp) != 0 {
		fmt.Printf("Run syslog server %s\n", *syslogIp)
		channel := make(syslog.LogPartsChannel)
		handler := syslog.NewChannelHandler(channel)

		server := syslog.NewServer()
		server.SetFormat(syslog.RFC3164)
		server.SetHandler(handler)
		server.ListenUDP(*syslogIp)
		server.Boot()
		go func(channel syslog.LogPartsChannel) {
			for logParts := range channel {
				msg := logParts["content"].(string)
				publisher.Push(msg)
				//fmt.Printf("%v\n", logParts)
			}
		}(channel)		
		server.Wait()		
		os.Exit(0)
	}
}
