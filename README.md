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

There are a number of [webdav clients](http://www.webdav.org/projects/).

## ChromeOS

For my specific use case, I am working with ChromeOS and there is a [WebDAV Storage Provider](https://chrome.google.com/webstore/detail/webdav-file-system/hmckflbfniicjijmdoffagjkpnjgbieh?hl=en).

## Linux

For Linux hosts, there is a package commonly `davfs2`, that provides a `mount.davfs` command.
See `mount.davfs(8)` man page for more information.

Basic example:
```bash
[vbatts@valse] {master} ~$ sudo mount.davfs https://bananaboat.usersys:9999/ ./x
Please enter the username to authenticate with server
https://bananaboat.usersys:9999/ or hit enter for none.
  Username: vbatts
Please enter the password to authenticate user vbatts with server
https://bananaboat.usersys:9999/ or hit enter for none.
  Password:  
mount.davfs: the server certificate is not trusted
  issuer:      Acme Co
  subject:     Acme Co
  identity:    localhost
  fingerprint: ce:19:a6:e7:0a:85:c2:01:fb:71:a6:bf:dd:56:3a:47:30:a8:7a:37
You only should accept this certificate, if you can
verify the fingerprint! The server might be faked
or there might be a man-in-the-middle-attack.
Accept certificate for this session? [y,N] y
[vbatts@valse] {master} ~$ ls x/file
x/file
[vbatts@valse] {master} ~$ cat !$
cat x/file
[vbatts@valse] {master} ~$ echo Howdy > !$
echo Howdy > x/file
[vbatts@valse] {master} ~$ cat x/file
Howdy
[vbatts@valse] {master} ~$ sudo umount ./x
/sbin/umount.davfs: waiting while mount.davfs (pid 8931) synchronizes the cache .. OK
```

