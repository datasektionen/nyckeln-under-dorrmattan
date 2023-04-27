# Nyckeln under dörrmattan
Like [login](https://github.com/datasektionen/login), but it always lets you in.

Can be used as a drop in replacement for login, but it requires no
configuration, and automatically lets everyone in as Ture Teknolog.

To build and run without docker and specifying a custom port:
```
go run . -port 1337
```

To build and run with docker and specifying a custom port:
```
docker build . -t login
docker run --rm --name login -p 1337:1337 login /login -port 1337
```

The default port is `10917`, so if you're happy with that, you don't have to
specify the `-port` flag.
