# mgraff-rc

# Test Driver

This command line tool is used to perform various operations to trigger
monitors to detect specific types of changes.

It can create files, delete them, and append to them.  It can run other
binaries, and will return the status code.  It can open TCP or UDP sockets
and transmit data.

All operations are logged to stdout as JSON.  In addition to the fields
specified in the requirements doc, additional debugging is provided.

At a minimum, each log message will contain a timestamp, called `ts`.

If an error is detected, such a a connection refused, a file cannot be created,
or other runtime error, a log message with `error` set will be generated, and
the application will panic.  This seems safest for now, as these opterations
are expected to succeed for the EDR agent to detect them, so any errors are
likely environmental or test framework errors.  In any case, the results will
not likely be useful, so failing fast seems best.

# Building

```
$ make
```

# Running

```
$ ./bin/testtool < commands.json
```

# commands.json

Each command is a single line of JSON, and will trigger a specific type of
action.

## Files

### Create File

The file must not already exist, and all directories must be created for the
provided path.

The file will be created with mode 0644.  It will be empty.

```
{
    "action": "CreateFile",
    "path": "/tmp/whateverPath"
}
```

### Modify File

This will write the provided `content` to the file at `path`.  The file must
already exist.

```
{
    "action": "ModifyFile",
    "path": "/tmp/whateverPath",
    "content": "something to append"
}
```

### Delete File

This will delete the file.  It if does not exist, this is not an error.

```
{
    "action": "DeleteFile",
    "path": "/tmp/whateverPath"
}
```

## Processes

### RunCommand

This will run the binary at `path`, passing the command-line arguments from the
`args` array.

If the binary does not exist, or is not executable, an error will be logged.
Otherwise, even if the exit code is non-zero, no error will be logged, and
the command will indicate success.

No timeout is provided; if the command does not exit, no logs will be created,
and the test code will hang.

The exit code will be returned on any successful execution.

```
{
    "action": "RunCommand",
    "path": "/bin/ls",
    "args": [ "/", "/tmp" ]
}
```

## Network

### NetworkWrite

Data can be written to any TCP or UDP destination, specified as either an IP
address or hostname.  If the address is provided, it may be either IPv4 or
IPv6, but IPv6 addresses must be included in `[]` format:  `[::1]`.  Link-local IPv6 addresses are also supported using the syntax `[fe80::1%lo0]`.

If using a DNS name, multiple addresses of both types may be returned.  To
force a connection using IPv4 or IPv6, "tcp4", "tcp6", "udp4", and "udp6" are
also supported as a protocol.

The bytes written indicate only that the data was inserted into the Kernel's
network stack buffers.  It does not guarantee that the data was successfully
sent.

```
{
    "action": "NetworkWrite",
    "protocol": "tcp",
    "host": "blog.flame.org",
    "port": 80,
    "data": "this is a test"
}
```
