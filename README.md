# srvdav

dangerously simple webdav server for a local filesystem

# Building

```bash
go get github.com/vbatts/srvdav
```

# Basic use

This daemon can serve up [WebDAV](https://en.wikipedia.org/wiki/WebDAV) for a local directory without any auth, nor encryption.
*DO NOT DO THIS*


# More proper use

Produce an x.509 certificate and accompanying key.
For development use case use can use the generator in golang's stdlib.

```bash
> go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -h
> go run $(go env GOROOT)/src/crypto/tls/generate_cert.go -host="localhost,example.com"
2016/09/22 09:46:19 written cert.pem
2016/09/22 09:46:19 written key.pem
```

Produce a password list for users.
The `htpasswd(1)` utility creates the password file nicely.

```bash
> htpasswd -bc srvdav.passwd vbatts topsecretpassword
```

Then launch `srvdav` with these credentials.

```bash
> mkdir -p ./test/
> srvdav -htpasswd ./srvdav.passwd -cert ./cert.pem -key ./key.pem
Serving HTTPS:// :9999
[...]
```

# Accompanying Clients

There are a number of webdav clients.
For my specific use case, I am working with ChromeOS and there is a [WebDAV Storage Provider](https://chrome.google.com/webstore/detail/webdav-file-system/hmckflbfniicjijmdoffagjkpnjgbieh?hl=en).

For Linux hosts, there is a package commonly `davfs2`, that provides a `mount.davfs` command.
See `mount.davfs(8)` man page for more information.

