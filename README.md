Deploy LiteFS on Fly.io
=======================

This repository provides a basic starting point for deploying LiteFS.


## Usage

### Deployment

First, create a `fly.toml` with the following contents:

```yml
[experimental]
  enable_consul = true

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    handlers = ["http"]
    port = 80
    force_https = true

  [[services.ports]]
    handlers = ["tls", "http"]
    port = "443"
```

Next, launch a new application:

```sh
fly launch
```

Choose "Y" when it asks to use the existing `fly.toml`. Choose an application
name and a region and then watch it deploy.


### Using the application

Once you have your application deployed, open a browser using your app's URL:

```
https://${APPNAME}.fly.dev/
```

