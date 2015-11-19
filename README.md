mediaserver
============

A web frontend to a media server.

# Configuration

Typical configuration options that you might have to change are in code/config/config-deploy.json.
"staticBaseDir" is the directory that you would like to expose.
"cacheBaseDir" should be an empty directory that you want to use for caching.

If you want to do any encoding (which you probably will have to do if you are serving any video files), then you will have to have ffmpeg and ffprobe installed.
Their paths can be set in the configuration file.

It is recommended that you only change configurations in code/config/config-deploy.json, but code/config/config-base.json contains all the configuration options.
Any options in code/config/config-deploy.json will override code/config/config-base.json.

If you are binding to a low port (like 80 or 443), then you may want to use something like setcap so you don't have to the entire server as root.

```
setcap 'cap_net_bind_service=+ep' /path/to/mediaserver
```

## Configuration Variables

| Variable              | Should Change?          | Default Value       | Description
|-----------------------|-------------------------|---------------------|------------
| staticBaseDir         | Yes                     | /srv/mediaserver    | This directory and all descendants will be served.
| cacheBaseDir          | Yes                     | /srv/cache          | The directory to use as a cache.
| encodingThreads       | Yes                     | 7                   | The number of threads to use while encoding.
| usersFile             | Maybe                   | config/users.json   | Path to the file that holds user information.
| port                  | Maybe                   | 80                  | The port to start the server on.
| showHiddenFiles       | Maybe                   | false               | Whether the server should show hidden files.
| ffmpegPath            | Maybe                   | /usr/bin/ffmpeg     | The path to ffmpeg.
| ffprobePath           | Maybe                   | /usr/bin/ffprobe    | The path to ffprobe.
| cacheUpperThresholdGB | Maybe                   | 50                  | When the cache grows larger then this, old entries will start to be removed.
| cacheLowerThresholdGB | Maybe                   | 40                  | When the cache is being compacted, entries will be removed until it is smaller than this.
| favicon               | Maybe                   | Webpage with camera | The hex for favicon.ico.
| clientBaseDir         | No                      | ../client           | Path to the directory to serve the client from.
| rawBaseURL            | No                      | raw                 | The URL prefix to use to serve static files.
| cacheBaseURL          | No                      | cache               | The URL prefix to use to serve cached files.
| clientBaseURL         | No                      | client              | The URL prefix to use to serve client files.
| apiVersion            | No                      | 0                   | The current API version.

# User Management

Users are managed through a json configuration file ("usersFile" in the configuration file).
Password hashes are stored in this file, so you should never share it.

There is currently no way to manage users online, but you can use the supplied manage-users utility to manage them.
Invoking the utility with no arguments will give you the usage.
