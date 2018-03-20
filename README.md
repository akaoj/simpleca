# simpleca

This tool allows you to easily generate and manage your private Certificate Authority. You can generate and sign root CAs, intermediate CAs and client keys.


# Usage

Let's say you want to generate a custom root CA, an intermediate CA and some client key pairs.

```
$ mkdir myca/ && cd myca/
$ simpleca generate root
Please provide the password for the file root/root.key:
Please repeat it:
Encrypted key generated in root/root.key
$ simpleca sign root
The file root/root.key is encrypted, please enter the password to unlock it:
root key signed, certificate available in root/root.crt
$ simpleca generate intermediate
Please provide the password for the file intermediates/intermediate.key:
Please repeat it:
Encrypted key generated in intermediates/intermediate.key
$ simpleca sign intermediate --with root
The file intermediates/intermediate.key is encrypted, please enter the password to unlock it:
The file root/root.key is encrypted, please enter the password to unlock it:
intermediate key signed, certificate available in intermediates/intermediate.crt
$ simpleca generate client
Please provide the password for the file clients/client.key:
Please repeat it:
Encrypted key generated in clients/client.key
$ simpleca sign client --with intermediate
The file clients/client.key is encrypted, please enter the password to unlock it:
The file intermediates/intermediate.key is encrypted, please enter the password to unlock it:
client key signed, certificate available in clients/client.crt
A full chain certificate file is also available at clients/client.fullchain.crt
$ ls clients/
client.crt  client.fullchain.crt  client.key  client.pub
```

You now simply have to configure your application to use these files.

Note that you can have multiple intermediates and clients:

```
$ simpleca generate intermediate --name intermediate01
Please provide the password for the file intermediates/intermediate01.key:
Please repeat it:
Encrypted key generated in intermediates/intermediate01.key
$ simpleca sign intermediate --name intermediate01 --with root
The file intermediates/intermediate01.key is encrypted, please enter the password to unlock it:
The file root/root.key is encrypted, please enter the password to unlock it:
intermediate01 key signed, certificate available in intermediates/intermediate01.crt
$ simpleca generate client --name web01.domain.com
Please provide the password for the file clients/web01.domain.com.key:
Please repeat it:
Encrypted key generated in clients/web01.domain.com.key
$ simpleca sign client --name web01.domain.com --with intermediate01
The file clients/web01.domain.com.key is encrypted, please enter the password to unlock it:
The file intermediates/intermediate01.key is encrypted, please enter the password to unlock it:
web01.domain.com key signed, certificate available in clients/web01.domain.com.crt
A full chain certificate file is also available at clients/web01.domain.com.fullchain.crt
$ ls clients
client.crt  client.fullchain.crt  client.key  client.pub  web01.domain.com.crt  web01.domain.com.fullchain.crt  web01.domain.com.key  web01.domain.com.pub
```

If you don't provide the `--name` flag, the default name will be used (`intermediate` for intermediate and `client` for client). Note that you can only have one root key pair and certificate.


# Configuration

When run for the first time, if no configuration is present, `simpleca` will generate one. You can then change the values as you like:

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

keyname = 'web01.domain.com'

httpd.socket = ssl.wrap_socket(
	httpd.socket,
	certfile='./clients/{}.fullchain.crt'.format(keyname),
	keyfile='./clients/{}.key'.format(keyname),
	server_side=True,
)
httpd.serve_forever()
```

Then you can try it with `curl`:

```
$ python server.py &
$ curl https://web01.domain.com:4443 -iv --cacert root/root.crt --resolve 'web01.domain.com:4443:127.0.0.1'
```
