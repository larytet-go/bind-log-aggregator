# Run the BIND

```
docker run --name bind -it --rm   --publish 5301:53/tcp --publish 5301:53/udp --publish 10000:10000/tcp   --volume $HOME/bind:/data --entrypoint "/usr/sbin/named" sameersbn/bind -c /data/named.conf -f -u bind
tail -F queries.log
dig @127.0.0.1 -p 5301 google.com
```

# Test the stdin 

```
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build bind-log-aggregator.go && sudo docker run --name bind -it --rm   --publish 5301:53/tcp --publish 5301:53/udp --publish 10000:10000/tcp    "/usr/sbin/named" sameersbn/bind -c /data/named.conf -f -u bind | ./bind-log-aggregator -logger stdout --logfile stdin
```

# Test the syslog 

```
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build bind-log-aggregator.go && sudo ./bind-log-aggregator -logger stdout --syslogip 127.0.0.1:514
logger -n 127.0.0.1 -P 514 "Hello, world"
```


# Setup rsyslog

Create a file inside the /etc/rsyslog.d folder called 10-rsyslog.conf
```
*.*   @remote.server:514
```
Restart the rsyslog service 
```
service rsyslog restart
```