# Timecop

## Building

```bash
make
```

## Running

NB: Currently the CLI blocks until the menubar app is quit (see [TODO](#TODO))

```bash
./timecop
```


# TODO

- add https://github.com/takama/daemon
  - attempted to run `menuet` inside a goroutine, but MacOS doesn't like that (needs to be on main thread)
  - integrating daemonisation so the CLI is both a daemon (which runs `menuet`) and a client to interact with the daemon

