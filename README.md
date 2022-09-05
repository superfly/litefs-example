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

Once you have your application deployed, open a browser using your app's URL.
Replace `${APPNAME}` with the name of your application on Fly.io.

```
https://${APPNAME}.fly.dev/
```

You can click the "Generate Record" button to create a new record on the primary.
If you make the request to a replica, it will automatically replay the request
on the primary.


### Deploying to multiple regions

Once you have your application running in one region, we can scale up easily by
using the `fly` command line interface. First, add additional regions:

```sh
fly regions add lhr syd
```

These 2 regions plus your original region you deployed to gives you 3 regions.
We'll need to increase our scale count to 3 to use all of them:

```sh
fly scale count 3
```

You should see the new regions start up in `fly logs` and `fly status`. Once
they're up and running, we can view their state by adding a `region` query
parameter to our URL:

```
https://${APPNAME}.fly.dev?region=lhr
```


### Watching replication in real-time

The application also includes a plain text response so we can use cURL to view
the data. We can use `watch` to continually run and display our results. Execute
the following commands in different terminal windows to watch the data in each
region:

```sh
watch -n 0.1 "curl -H 'Accept: text/plain' https://${APPNAME}.fly.dev?region=lhr"
```

```sh
watch -n 0.1 "curl -H 'Accept: text/plain' https://${APPNAME}.fly.dev?region=syd"
```

Now when you click the "Generate Record", you should see the data show up in
your other regions almost immediately. The replication time is mostly dependent
on the speed of light so you may see delays of around 250ms if the data needs to
travel half way around the world.

