# mgraff-rc

# Test Driver

This command line tool is used to perform various operations to trigger
monitors to detect specific types of changes.

It can create files, delete them, and append to them.  It can run other
binaries, and will return the status code.  It can open TCP or UDP sockets
and transmit data.

All operations are logged to stdout as JSON.  In addition to the fields
specified in the requirements doc, additional debugging is provided.

At a minimum, each log message will contain a timestamp, called `ts`,
`processID`, `processName`, `processArgs`, and `processUsername`.

* `ts` is the Unix timestamp in Epoch seconds, and may be a floating point value.
* `processID` is the PID of the running process.
* `processName` is the full path to the running process's executable, aka the "command".
* `processArgs` are any command line arguments passed to the process.
* `processUsername` is the username running the process.

If an error is detected, such a a connection refused, a file cannot be created,
or other runtime error, a log message with `error` set will be generated, and
the application will panic.  This seems safest for now, as these opterations
are expected to succeed for the EDR agent to detect them, so any errors are
likely environmental or test framework errors.  In any case, the results will
not likely be useful, so failing fast seems best.

For any action, if `error` is not present in the JSON log, and `action` is,
it indicates a successful operation.

If `action` is not set, it is a normal log message and can be ignored for EDR
activity comparison.

# Building

```
$ make
```

# Running

```
$ ./bin/testtool < commands.json
```

# commands.json

This contains a sequence of JSON action descriptions.  Each is an object, and a blank line indicates the end of the JSON object.

See `test.json` for an example.

# Actions

## Files

### Create File

The file must not already exist, and all directories must be created for the
provided path.

The file will be created with mode 0644.  It will be empty.

```
// Request
{
    "action": "CreateFile",
    "path": "/tmp/whateverPath"
}
```

```
// Response
{
    "ts":1655242678.3220391,
    "processID":21453,
    "processUsername":"explorer",
    "processName":"/tmp/bin/testtool",
    "processArguments":[],
    "action":"CreateFile",
    "fileAction":"create",
    "path":"/tmp/foo"
}
```

### Modify File

This will append the provided `content` to the file at `path`.  The file must
already exist.

```
// Request
{
    "action": "ModifyFile",
    "path": "/tmp/whateverPath",
    "content": "something to append"
}
```

```
// Response
{
    "ts": 1655242678.32215,
    "processID": 21453,
    "processUsername": "explorer",
    "processName":"/tmp/bin/testtool",
    "processArguments": [],
    "action": "ModifyFile",
    "fileAction": "modify",
    "path": "/tmp/foo"
}
```

### Delete File

This will delete the file.  It if does not exist, this is not an error.

```
// Request
{
    "action": "DeleteFile",
    "path": "/tmp/whateverPath"
}
```

```
// Response
{
    "ts": 1655242678.3224661,
    "processID": 21453,
    "processUsername": "explorer",
    "processName":"/tmp/bin/testtool",
    "processArguments": [],
    "action": "DeleteFile",
    "fileAction": "delete",
    "path": "/tmp/foo"
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
// Request
{
    "action": "RunCommand",
    "path": "/bin/ls",
    "args": [ "/", "/tmp" ]
}
```

```
// Response
{
    "ts": 1655242678.325207,
    "processID": 21453,
    "processUsername": "explorer",
    "processName":"/tmp/bin/testtool",
    "processArguments": [],
    "action": "RunCommand",
    "cmdPath": "/bin/ls",
    "cmdArgs": [
        "/foo"
    ],
    "cmdPID": 21454,
    "cmdExitStatus": 1
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
// Request
{
    "action": "NetworkWrite",
    "protocol": "tcp",
    "host": "blog.flame.org",
    "port": 80,
    "data": "this is a test"
}
```

```
// Response
{
    "ts": 1655242678.445803,
    "processID": 21453,
    "processUsername": "explorer",
    "processName":"/tmp/bin/testtool",
    "processArguments": [],
    "action": "NetworkWrite",
    "host": "blog.flame.org",
    "port": 80,
    "protocol": "tcp",
    "localAddress": "10.42.0.14",
    "localPort": 52678,
    "localZone": "",
    "nWritten": 14
}
```
