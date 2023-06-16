LiteFS Example Application on Fly.io
====================================

## Prerequisites

First, you'll need to install [`flyctl`][] and then [sign up for a free account][signup].

[`flyctl`]: https://fly.io/docs/hands-on/install-flyctl/
[signup]: https://fly.io/docs/hands-on/sign-up/


## Usage

### Creating the application

First, we'll need to create our application on [Fly.io](https://fly.io). This
will prompt you for a name or you can autogenerate one for this example.
Remember this name as we'll need it later.

```sh
fly launch --region ord --no-deploy
```

The launch command will create a `fly.toml` file and set the primary region to
Chicago (`ord`) but will not launch the app.


### Creating a volume

Next, we need to set up a persistent volume in our primary region so that our
data is not lost between restarts.

```sh
fly volumes create -r ord --size 1 litefs
```

And add a mount to this volume in your `fly.toml` file:

```toml
[mounts]
  source = "litefs"
  destination = "/var/lib/litefs"
```

### Setting up Consul

LiteFS uses [Consul](https://consul.io) for its distributed lease. You can find
instructions for using Fly.io's free multi-tenant Consul in the
[Lease Management][] section of the Getting Started guide.

[Lease Management]: https://fly.io/docs/litefs/getting-started/#lease-configuration


### Launching your app

The next step is to launch & deploy your app with the following command:

```sh
fly deploy
```

The application should build and deploy and you should see it up in running
after a minute or so. You can go to `https://$APPNAME.fly.dev/` and see your
application running live.

The application is a simple interface for generating fake records. It's just
for illustrating how LiteFS can easily replicate your data between nodes.

When you click the _"Generate Record"_ button, it will create that row in a
local SQLite database that is running on a LiteFS file system. Any other node
running LiteFS will automatically get those updates and apply them to their
local copy of the database. That lets every node keep an exact copy of the same
database.


### Launching more regions

This example application is configured to run as a primary only in the
`PRIMARY_REGION` (which is Chicago). It's best practice to run two or more
instances in the primary region and then you can add instances in additional
regions to reduce latency for your users.

You can clone the configuration of the machine to other regions by using the
`fly m clone` command. The `--select` flag lets you choose from a list of
existing machines to clone.

```sh
# Make a second instance in your primary region.
fly m clone --select --region ord

# Make additional instances in regions around the world (London, Sydney, etc).
fly m clone --select --region lhr
fly m clone --select --region syd
```
