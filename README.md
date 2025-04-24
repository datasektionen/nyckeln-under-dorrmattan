# Nyckeln under dörrmattan

Mock version of [login](https://github.com/datasektionen/login),
[pls](https://github.com/datasektionen/pls), and [sso](https://github.com/datasektionen/sso).

The login part can be used as a drop in replacement, but it requires no
configuration, and automatically lets everyone in as Ture Teknolog. You
can also write someones kth id on stdin, which will make future login
requests be for that user. Data is then fetched from hodis.

The pls part only implements a subset of the pls API. Feel free to extend
it if you need more. If that isn't good enough, use the pls in production
and log in as someone with enough privileges, such as the current kassör or d-sys.

<details>
<summary>Pls API</summary>
<br>

* `GET /api/user/:id`, returns all map of groups with a list of permissions for a user
* `GET /api/user/:id/:group`, returns a list all group permissions for a user
* `GET /api/user/:id/:group/:permission`, returns true or false if a user has the permission

</details>

Lastly, the sso part is a simple OpenID Connect (oidc) server which behaves the
same way that sso does. You define your configuration and users in a yaml file,
and then configure your favorite oidc client to use `http://localhost:{sso-port}/.well-known/openid-configuration`,
and it should just work. When logging in, you will need to enter the username/kth-id of
someone defined in your yaml config. Similarly to sso, also supports `pls_*` scopes.

## Configuration

You can configure the following flags:

* `pls-port`: Port for the pls service. Defaults to 7001.
* `login-port`: Port for the login service. Defaults to 7002.
* `sso-port`: Port for the sso service. Defaults to 7003.
* `hodis-url`: URL to the hodis instance. Defaults to `https://hodis.datasektionen.se`.
* `kth-id`: Username to use for login. Defaults to `KTH_ID` environment variable, or `turetek` if not set.
* `config-file`: Path to a yaml config file. Defaults to `config.yaml`.

The yaml config file is used for SSO (oidc) and pls configuration. For example:

```yaml
# OpenID Client configuration
clients:
  - id: "client-id" # use this in your oidc client
    secret: "client-secret" # use this in your oidc client
    redirect_uris: # URIs of your oidc client
      [
        "http://localhost:4000/oidcc/callback",
        "http://localhost:4000/oidcc/authorize",
      ]

# Users to log in via sso, and their pls permissions
users:
  - ug_kth_id: some-id
    kth_id: turetek
    email: turetek@kth.se
    first_name: Ture
    family_name: Teknolog
    pls_permissions:
      sso: [admin]
      calypso:
        - drek
        - dfunk
```


To build and run without docker and specifying custom ports:
```
go run . -login-port 1337 -pls-port 1338 -sso-port 1339 -config-file my.config.yaml
```

To build and run with docker using the default ports:
```
docker build . -t nyckeln
docker run -it --rm --name nyckeln \
    -v /path/to/your/config.yaml:/config.yaml \
    -p 7001:7001 -p 7002:7002 -p 7003:7003 nyckeln
```

The container is also published as a container at
ghcr.io/datasektionen/nyckeln-under-dorrmattan, so you can also run it as

```
docker run -it --rm --name nyckeln \
    -v /path/to/your/config.yaml:/config.yaml \
    -p 701:7001 -p 7002:7002 -p 7003:7003 ghcr.io/datasektionen/nyckeln-under-dorrmattan
```

Without even having to clone this repository. You can also add it to your dev
`docker-compose.yaml` file.

```yaml
services:
  nyckeln:
    image: ghcr.io/datasektionen/nyckeln-under-dorrmattan
    configs:
      - source: nyckeln.yaml
        target: /config.yaml
    ports:
      - 7001:7001
      - 7002:7002
      - 7003:7003

configs:
  nyckeln.yaml:
    content: |
      clients:
        - id: "client-id"
          secret: "client-secret"
          redirect_uris:
            - "http://localhost:4000/oidcc/callback"
            - "http://localhost:4000/oidcc/authorize"

      users:
        - kth_id: turetek
          email: turetek@kth.se
          first_name: Ture
          family_name: Teknolog
          pls_permissions:
            sso:
              - fippel
```
