Deploy LiteFS on Fly.io
=======================

This repository provides a basic starting point for deploying LiteFS. To run it,
simply run the following commands:

```sh
fly launch
```

Choose "Y" when it asks to use the existing `fly.toml`. Choose an application
name and a region and then watch it deploy.

Once deployed, you can log in using:

```sh
fly ssh console
```

Inside the Firecracker VM, you can execute SQLite commands against the
LiteFS mounted file system at `/data`:

```sh
sqlite3 /data/db
sqlite3> CREATE TABLE widgets (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);
```
