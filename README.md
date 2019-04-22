Run the BIND
```
docker run --name bind -it --rm   --publish 5301:53/tcp --publish 5301:53/udp --publish 10000:10000/tcp   --volume $HOME/bind:/data --entrypoint "/usr/sbin/named" sameersbn/bind -c /data/named.conf -f -u bind
tail -F queries.log
dig @127.0.0.1 -p 5301 google.com
```


Test the syslog 
```
logger -n 127.0.0.1 -P 514 "Hello, world"
```
