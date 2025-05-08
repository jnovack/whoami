# whoami

Tiny Go webserver to print HTTP request vars to output.

![whoami logo](docs/logo.svg)

```sh
$ docker run -d -P --name whoami jnovack/whoami
49e9ee7948300c7f310ab9f6e405c658f8433cc418e6f94e36e1171524209ed8

$ docker port whoami
80/tcp -> 0.0.0.0:32769

$ curl http://localhost:32769
IP: 172.17.0.2
ENV: PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV: HOSTNAME=ce9538e7bfe2
ENV: HOME=/nonexistent
BUILD_VERSION: v1.4.2-2-ga35afdf
BUILD_COMMIT: a35afdfa7fe43c30c562c90eca172f2eead7a468
BUILD_RFC3339: 2025-04-18T14:47:34+00:00
TIMESTAMP: 2025-04-18T14:48:18.348892168Z
PROTOCOL: HTTP/1.1
GET / HTTP/1.1
Host: localhost:32769
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:137.0) Gecko/20100101 Firefox/137.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Encoding: gzip, deflate, br, zstd
Accept-Language: en-US,en;q=0.5
Cache-Control: must-validate
Connection: keep-alive
Hostname: ce9538e7bfe2
Priority: u=0, i
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: cross-site
Upgrade-Insecure-Requests: 1
```

Forked from [emilevauge/whoami](https://github.com/emilevauge/whoamI).
