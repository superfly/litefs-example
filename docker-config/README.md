LiteFS Example Application with Docker
======================================

## Prerequisites

You need to [install Docker][] (and docker-compose, which is included in all
recent releases of Docker).

[install Docker]: https://docs.docker.com/engine/install/

## Usage

### Running the app

You can run the application with this command (from the top-level directory
of the repo):

```bash
docker-compose up
```

### Using the app locally

The app is configured with the following containers:

* `primary` - this is the application with LiteFS running as primary
(all writes happen on this node)
* `replica1` - this is the application running with LiteFS as a replica
* `nginx` - nginx configured as a load balancer, and also configured to
route all `POST` requests to the primary node

Each of these is available to test out locally:

* http://localhost:8080 - the load balancer
* http://localhost:8081 - the primary
* http://localhost:8082 - the replica

You can test generating records on each node (it will fail if attempted on
the replica).