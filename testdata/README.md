# SETUP

Use: https://github.com/golang/go/blob/master/src/crypto/tls/generate_cert.go

```
% wget https://raw.githubusercontent.com/golang/go/master/src/crypto/tls/generate_cert.go
% go run generate_cert.go --host="localhost"
```

```
r.RunTLS("", "./testdata/cert.pem", "./testdata/key.pem")
```

---

Others:

```
% openssl req -newkey rsa:4096 -nodes -keyout server.key -out server.csr
% openssl x509 -signkey server.key -in server.csr -req -days 365 -out server.crt
```

```
r.RunTLS("", "./testdata/server.crt", "./testdata/server.key")
```

# Chrome Browser

```
thisisunsafe
```
