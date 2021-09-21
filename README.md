## Fleetlock client

[![Go Reference](https://pkg.go.dev/badge/github.com/flatcar-linux/fleetlock.svg)](https://pkg.go.dev/github.com/flatcar-linux/fleetlock)
[![Go](https://github.com/flatcar-linux/fleetlock/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/flatcar-linux/fleetlock/actions/workflows/go.yml)

Go implementation of `FleetLock` protocol.

### Example

It's possible to run a local `FleetLock` server backed by `etcd` using [`airlock`](https://github.com/coreos/airlock/):

```shell
$ cd examples
$ docker-compose up -d
```
There is only one slot for the `controllers` group:
```
$ fleetlockctl --url http://127.0.0.1:3333 --group controllers --id $(cat /etc/machine-id) recursive-lock
```

On the `airlock` side:
```
DEBU[0117] processing client pre-reboot request          group=controllers id=41987...
DEBU[0117] givin green-flag to pre-reboot request        group=controllers id=41987...
```

It's not possible to get another lock, since there is only one slot:
```
$ fleetlockctl --url http://127.0.0.1:3333 --group controllers --id another-machine recursive-lock
Error: locking: fleetlock error: all 1 semaphore slots currently locked (failed_lock)
```

It's possible to release the lock:
```
./fleetlockctl --url http://127.0.0.1:3333 --group controllers --id $(cat /etc/machine-id) unlock-if-held
```

`another-machine` can now request a lock:
```
$ fleetlockctl --url http://127.0.0.1:3333 --group controllers --id another-machine recursive-lock
```

On the `airlock` side:
```
DEBU[0366] got pre-reboot request
DEBU[0366] processing client pre-reboot request          group=controllers id=another-machine
DEBU[0367] givin green-flag to pre-reboot request        group=controllers id=another-machine
```

### Build

requirements:
  * `go` in the path

```
$ make
```

ref: https://coreos.github.io/zincati/development/fleetlock/protocol/
