mediaserver
============

A web frontend to a media server.

If you want the version that does encoding, use the `encode` branch.
To simplify, encoding and cache work has been dropped.

# Configuration

You can see a full sample configuration file in code/config/config-sample.json.

If you are binding to a low port (like 80 or 443), then you may want to use something like setcap so you don't have to the entire server as root.

```
setcap 'cap_net_bind_service=+ep' /path/to/mediaserver
```

## Configuration Variables

| Variable              | Should Change? | Default Value           | Description
|-----------------------|----------------|-------------------------|------------
| staticBaseDir         | Yes            | /srv/mediaserver        | This directory and all descendants will be served.
| usersFile             | Maybe          | config/users.json       | Path to the file that holds user information.
| prodConfig            | Maybe          | config/config-prod.json | Path to the production configuration file. When the --prod flag is supplied, this configuration will be applied after the base configuration.
| httpPort              | Maybe          | 80                      | The port to start the server on when using http.
| useSSL                | Maybe          | false                   | Whether or not yo use SSL.
| httpsPort             | Maybe          | 443                     | The port to start the server on when using https.
| forwardHttp           | Maybe          | true                    | If using https and this option is active, then the http port will be forwarded to the https port.
| sslCertFile           | Maybe          | config/domain.cert      | The cert file to use for SSL.
| sslKeyFile            | Maybe          | config/domain.key       | The key file to use for SSL.
| showHiddenFiles       | Maybe          | false                   | Whether the server should show hidden files.
| favicon               | Maybe          | Webpage with camera     | The hex for favicon.ico.
| clientBaseDir         | No             | ../client               | Path to the directory to serve the client from.
| rawBaseURL            | No             | raw                     | The URL prefix to use to serve static files.
| clientBaseURL         | No             | client                  | The URL prefix to use to serve client files.
| apiVersion            | No             | 0                       | The current API version.

# User Management

Users are managed through a json configuration file ("usersFile" in the configuration file).
Password hashes are stored in this file, so you should never share it.

There is currently no way to manage users online, but you can use the supplied manage-users utility to manage them.
Invoking the utility with no arguments will give you the usage.

## User Management Examples

Adding a user:
```
code/bin/manage-users add config/users.json
```
Add a user to a non-existant file will just create a new file populated with that user.

List the users:
```
code/bin/manage-users ls config/users.json
```
