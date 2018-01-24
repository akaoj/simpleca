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
$ simpleca sign client --with intermediate01
$ simpleca show cert client > web01.domain.com.crt
$ simpleca show priv client > web01.domain.com.key
$ simpleca show pub client > web01.domain.com.pub
```

You now simply have to configure your application to use these files.

Note that you can have multiple intermediates and clients:

```
$ simpleca generate intermediate --name intermediate01
$ simpleca sign intermediate --name intermediate01 --with root
$ simpleca generate client --name web01.domain.com
$ simpleca sign client --name web01.domain.com --with intermediate01
$ simpleca show cert client --name web01.domain.com
```

If you don't provide the `--name` flag, the default name will be used (`intermediate` for intermediate and `client` for client). Note that you can only have one root key pair / certificate.
