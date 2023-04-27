# Nyckeln under dörrmattan
Like [login](https://github.com/datasektionen/login) and
[pls](https://github.com/datasektionen/pls), but it always lets you in and
thinks you can do anything.

The login part can be used as a drop in replacement, but it requires no
configuration, and automatically lets everyone in as Ture Teknolog.

The pls part only works when validating that a user/mandate/token has a given
permission. It can not list the given permissions for someone.

To build and run without docker and specifying custom ports:
```
go run . -login-port 1337 -pls-port 1338
```

To build and run with docker using the default ports:
```
docker build . -t nyckeln
docker run --rm --name nyckeln \
    -p 7001:7001 -p 7002:7002 nyckeln
```
