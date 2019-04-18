Run the BIND
```
docker run --name bind -it --rm   --publish 5301:53/tcp --publish 5301:53/udp --publish 10000:10000/tcp   --volume $HOME/bind:/data --entrypoint "/usr/sbin/named" sameersbn/bind:9.11-20190315 -c /data/named.conf -f -u bind
```
