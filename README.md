# simpleca

This tool allows you to easily generate and manage your private Certificate Authority. You can generate and sign root CAs, intermediate CAs and client keys.


# Usage

Let's say you want to generate a custom root CA, an intermediate CA and some client key pairs.

```
$ mkdir myca/ && cd myca/
$ simpleca generate root
$ simpleca sign root
$ simpleca generate intermediate
$ simpleca sign intermediate --with root
$ simpleca generate client
$ simpleca sign client --with intermediate
$ ls clients/
client.crt  client.key  client.pub
```

You now simply have to configure your application to use these files.

Note that you can have multiple intermediates and clients:

```
$ simpleca generate intermediate --name intermediate01
$ simpleca sign intermediate --name intermediate01 --with root
$ simpleca generate client --name web01.domain.com
$ simpleca sign client --name web01.domain.com --with intermediate01
$ ls clients
client.crt  client.key  client.pub  web01.domain.com.crt  web01.domain.com.key  web01.domain.com.pub
```

If you don't provide the `--name` flag, the default name will be used (`intermediate` for intermediate and `client` for client). Note that you can only have one root key pair and certificate.


# Configuration

When run for the first time, if no configuration is present, `simpleca` will generate one. You can the change the values as you like:

- CertificateDuration: specify the duration of signed certificates **in months**
- Organization: the name of your organization
- Country: your country
- Locality: your city

Note that these informations are **only** used for the certificates. They are **not** and **never will be** sent to some strange remote server and are **not** used for statistics purposes.


# Test it

Spawn a simple HTTPS server:

`server.py`:

```python
import BaseHTTPServer, SimpleHTTPServer
import ssl

httpd = BaseHTTPServer.HTTPServer(
	('localhost', 4443),
	SimpleHTTPServer.SimpleHTTPRequestHandler
)
httpd.socket = ssl.wrap_socket(
	httpd.socket,
	certfile='./clients/web01.domain.com.fullchain.crt',
	keyfile='./clients/web01.domain.com.key',
	server_side=True,
)
httpd.serve_forever()
```

Then you can try it with `curl`:

```
$ cat clients/web01.domain.com.crt intermediates/intermediate01.crt > clients/web01.domain.com.fullchain.crt
$ python server.py &
$ curl https://web01.domain.com:4443 -iv --cacert root/root.crt --resolve 'web01.domain.com:4443:127.0.0.1'
```
