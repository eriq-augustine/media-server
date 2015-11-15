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

# User Management

Users are managed through a json configuration file ("usersFile" in the configuration file).
Password hashes are stored in this file, so you should never share it.

There is currently no way to manage users online, but you can use the supplied manage-users utility to manage them.
Invoking the utility with no arguments will give you the usage.
