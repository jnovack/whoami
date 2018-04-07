# whoami

Tiny Go webserver that print HTTP request to output.

```sh
$ docker run -d -P --name whoami jnovack/whoami
49e9ee7948300c7f310ab9f6e405c658f8433cc418e6f94e36e1171524209ed8

$ docker port whoami
80/tcp -> 0.0.0.0:32769

$ curl http://0.0.0.0:32769
Hostname: 49e9ee794830
IP: 192.168.203.6
ENV: PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV: HOSTNAME=49e9ee794830
ENV: HOME=/
VERSION: v1.2.0
COMMIT: 4ded474
BUILD_DATE: 2018-04-06
BUILD_TIME: 20:46:43-0400
GET / HTTP/1.1
Host: localhost:32769
User-Agent: curl/7.54.0
Accept: */*
```

Forked from [emilevauge/whoami](https://github.com/emilevauge/whoamI).
