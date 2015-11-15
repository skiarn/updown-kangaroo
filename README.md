# updown-kangaroo

This application is to be used to notify when systems are down. Survailance software or proxy servers should copy 503 requests and send to this application. The application listens for http 503 errors on path: ```localhost:8080/in```.

updown-kangaroo will distribute 503 errors over websockets to subscribers using following json data format:
```{"host":"","path":"/in"}```.
Subscribers can then notify their users about unavailable systems.

## Run application
```
nohup ./updown-kangaroo -addr=:8080 > server.log 2>&1 &
```

