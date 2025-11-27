# Nyckeln under dörrmattan

Mock version of [login](https://github.com/datasektionen/login),
[pls](https://github.com/datasektionen/pls), [hive](https://github.com/datasektionen/hive),
[sso](https://github.com/datasektionen/sso), and [ldap-proxy](https://github.com/datasektionen/ldap-proxy).

The login part can be used as a drop in replacement, but it requires no
configuration, and automatically lets everyone in as Ture Teknolog. You
can also write someones kth id on stdin, which will make future login
requests be for that user. Data is then fetched from hodis.

The pls part only implements a subset of the pls API. Feel free to extend
it if you need more. If that isn't good enough, use the pls in production
and log in as someone with enough privileges, such as the current kassör or d-sys. The
same applies to Hive.

<details>
<summary>Pls API</summary>
<br>

- `GET /api/user/:id`, returns all map of groups with a list of permissions for a user
- `GET /api/user/:id/:group`, returns a list all group permissions for a user
- `GET /api/user/:id/:group/:permission`, returns true or false if a user has the permission

</details>

<details>
<summary>Hive API</summary>
<br>

As in real Hive, you need to supply an `Authentication` header with `Bearer <token>` to access the API. Any token will be accepted.

- `GET /api/v1/user/:id/permissions`, returns a list of hive permissions for a user
- `GET /api/v1/user/:id/permission/:perm_id`, returns true or false if a user has the permmission
- `GET /api/v1/user/:id/permissions/:perm_id/scopes`, returns a list of scopes for a specific permission
- `GET /api/v1/user/:id/permissions/:perm_id/scope/:scope`, returns true or false if a user has a specific scope

- `GET /api/v1/token/:secret/permissions`, returns a list of hive permissions for a token
- `GET /api/v1/token/:secret/permission/:perm_id`, returns a list of hive permissions for a token
- `GET /api/v1/token/:secret/permission/:perm_id/scopes`, returns a list of scopes for a specific permission
- `GET /api/v1/token/:secret/permission/:perm_id/scope/:scope`, returns true or false if the token has a specific scope

- `GET /api/v1/tagged/:tag_id/groups`, returns a list of groups that have that tag
- `GET /api/v1/tagged/:tag_id/memberships/:username`, returns all the groups that the user is part of and has the tag
- `GET /api/v1/tagged/:tag_id/users`, returns a list of all users who have the tag
- `GET /api/v1/group/:group_domain/:group_id/members`, returns a list of all members in group
</details>

<details>
<summary>SSO API</summary>
<br>

Lastly, the sso part is a simple OpenID Connect (oidc) server which behaves the
same way that sso does. You define your configuration and users in a yaml file,
and then configure your favorite oidc client to use `http://localhost:{sso-port}/.well-known/openid-configuration`,
and it should just work. When logging in, you will need to enter the username/kth-id of
someone defined in your yaml config. If allow guest is turned on it can also return users form the ldap part of the yaml.
Similarly to sso, also supports `pls_*`, `permissions`, and `picture` scopes.

- `GET /api/users`, takes a list of users kthid (using repeated u query parameters) and a format query parameter
    and returns user information based on the format. It also optionaly provides the users profile picture when supplied
    with the picture query parameter.
- `GET /api/search`, takes a query, offset, limit, and year parameters. It can also optionaly supply the users profile picture.
</details>

<details>
<summary>ldap-proxy API</summary>
<br>

The ldap-porxy part mocks the systems only endpoint but without the ability to search for ug_kth_id with a simple config interface in the yaml file.

- `GET /user`, takes a kthid as a query parameter and returns basic information about the user if it exists.
</details>

## Configuration

You can configure the following flags:

- `pls-port`: Port for the pls service. Defaults to 7001.
- `login-port`: Port for the login service. Defaults to 7002.
- `sso-port`: Port for the sso service. Defaults to 7003.
- `hive-port`: Port for the hive service. Defaults to 7004.
- `ldap-port`: Port for the kthldap service. Defaults to 7005.
- `hodis-url`: URL to the hodis instance. Defaults to `https://hodis.datasektionen.se`.
- `kth-id`: Username to use for login. Defaults to `KTH_ID` environment variable, or `turetek` if not set.
- `config-file`: Path to a yaml config file. Defaults to `config.yaml`.

The yaml config file is used for SSO (oidc) and pls configuration. For example:

```yaml
# OpenID Client configuration
clients:
  - id: "client-id" # use this in your oidc client
    secret: "client-secret" # use this in your oidc client
    allows_guests: false
    redirect_uris: # URIs of your oidc client
      - http://localhost:4000/oidcc/callback

# Users to log in via sso, and their pls permissions
users:
  - ug_kth_id: some-id
    kth_id: turetek
    email: turetek@kth.se
    first_name: Ture
    family_name: Teknolog
    picture: https://thskth.se/wp-content/uploads/2023/06/ths-logo-bla-01-300x294-1.png
    thumbnail: https://thskth.se/wp-content/uploads/2019/12/ths-logo-svart-01.png
    pls_permissions:
      sso:
        - admin
      calypso:
        - drek
        - dfunk
    hive_tags:
      - id: personal_email
        content: turetek@gmail.com

# Hive setup including tokens and groups with permissions and tags
hive:
  tokens:
    - secret: aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
      permissions:
        - id: admin
          scope: null
  groups:
    - name: Systemansvarig
      id: d-sys
      domain: example.com
      members:
        - turetek
      tags:
        - id: author-pseudonym
          content: D-Sys
        - id: mandate
      permissions:
        - id: admin
          scope: null

# Ldap proxy setup
ldap:
  - ug_kth_id: some-other-id
    kth_id: marmas
    first_name: Markus
    family_name: Maskinare
```

To build and run without docker and specifying custom ports:

```
go run . -login-port 1337 -pls-port 1338 -sso-port 1339 -hive-port 1340 -config-file my.config.yaml
```

To build and run with docker using the default ports:

```
docker build . -t nyckeln
docker run -it --rm --name nyckeln \
    -v /path/to/your/config.yaml:/config.yaml \
    -p 7001:7001 -p 7002:7002 -p 7003:7003 -p 7004:7004 nyckeln
```

The container is also published as a container at
ghcr.io/datasektionen/nyckeln-under-dorrmattan, so you can also run it as

```
docker run -it --rm --name nyckeln \
    -v /path/to/your/config.yaml:/config.yaml \
    -p 701:7001 -p 7002:7002 -p 7003:7003 -p 7004:7004 ghcr.io/datasektionen/nyckeln-under-dorrmattan
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
      - 7004:7004

configs:
  nyckeln.yaml:
    content: |
      clients:
        - id: "client-id" # use this in your oidc client
          secret: "client-secret" # use this in your oidc client
          allows_guests: false
          redirect_uris: # URIs of your oidc client
            - http://localhost:4000/oidcc/callback

      users:
        - ug_kth_id: some-id
          kth_id: turetek
          email: turetek@kth.se
          first_name: Ture
          family_name: Teknolog
          picture: https://thskth.se/wp-content/uploads/2023/06/ths-logo-bla-01-300x294-1.png
          thumbnail: https://thskth.se/wp-content/uploads/2019/12/ths-logo-svart-01.png
          pls_permissions:
            sso:
              - admin
            calypso:
              - drek
              - dfunk
          hive_tags:
            - id: personal_email
              content: turetek@gmail.com

      hive:
        tokens:
          - secret: aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa
            permissions:
              - id: admin
                scope: null
        groups:
          - name: Systemansvarig
            id: d-sys
            domain: example.com
            members:
              - turetek
            tags:
              - id: author-pseudonym
                content: D-Sys
              - id: mandate
            permissions:
              - id: admin
                scope: null

      ldap:
        - ug_kth_id: some-other-id
          kth_id: marmas
          first_name: Markus
          family_name: Maskinare
```
