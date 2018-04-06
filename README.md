# whoami

Tiny Go webserver that print HTTP request to output.

```sh
$ docker run -d -P --name iamfoo jnovack/whoami
$ docker inspect --format '{{ .NetworkSettings.Ports }}' iamfoo
map[80/tcp:[{0.0.0.0 32769}]]
$ curl "http://0.0.0.0:32769"
Hostname :  6e0030e67d6a
IP :  127.0.0.1
IP :  ::1
IP :  172.17.0.27
IP :  fe80::42:acff:fe11:1b
GET / HTTP/1.1
Host: 0.0.0.0:32769
User-Agent: curl/7.35.0
Accept: */*
```

Forked from emilevauge/whoami.
