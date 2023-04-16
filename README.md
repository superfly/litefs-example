LiteFS Example Application
==========================

This repository is an example of a toy application running on LiteFS on Fly.io.

## Prerequisites

First, you'll need to install [`flyctl`](/docs/hands-on/install-flyctl/) and
then [sign up for a free account](https://fly.io/docs/hands-on/sign-up/).


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

