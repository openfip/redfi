<img src="https://raw.githubusercontent.com/redfi/redfi/master/static/redfi@2x.png" width="400px">

RedFI acts as a proxy between the client and Redis with the capability
of injecting faults on the fly, based on the rules given by you.

## Features
- Simple to use. It is just a binary that you execute.
- Transparent to the client.
- Flexibility to inject failures on Redis Commands you define.
- Limit failure injection to a percentage of the commands.
- Limit failure injection to certain clients only.
- Inject latency, drop connections, return empty responses.

## Why?
We believe that to gain trust in your system,
you must test it against various scenarios of failures, before deploying it to production.  
These tests should also be part of your CI, to catch code changes that might introduce regressions.  
This would help you understand the limits of your system, and how resilient it is against Redis' failures.

## How it Works
RedFI is a proxy that sets between the client and the actual Redis server. On every incoming command from the client, it checks the list of failure rules provided by you, and then it applies the first rule that matches the request.

## Usage
First, download the [latest binary](https://github.com/redfi/redfi/releases/tag/v0.1).

After that, execute the binary we downloaded earlier
```bash
$ ./redfi -addr 127.0.0.1:8083 -redis 127.0.0.1:6379

RedFI is listening on 127.0.0.1:8083
Don't forget to point your client to that address.

RedFI Controller is listening on :6380
```
- **addr**: is the address on which the proxy listens on for new connections.
- **redis**: address of the actual Redis server to proxy commands/connections to.
- **plan**: path to the json file that contains the rules/scenarios for fault injection.

An example of controlling the proxy rules
[![asciicast](https://asciinema.org/a/CJyxhZXfk7aaGJ9qfU865Gcy6.png)](https://asciinema.org/a/CJyxhZXfk7aaGJ9qfU865Gcy6)

## Config Commands

Point redis-cli to RedFI Controller port
```bash
$ redi-cli -p 6380
```

### RULEADD rule_name option=value [option=value]
Adds a rule to the engine of RedFI

### Faults Options
- `delay`: adds a delay, the value unit is milliseconds. Example: `delay=200`.
- `return_empty`: returns an empty response. Example: `return_empty=true`.
- `return_err`: returns an error response. Example: `return_err=ERR bad command.`
- `drop`: drop connections, the value is a boolean. Example: `drop=true`.

### Blast limiters
They limit the radius of the blast
- `command`: only apply failure on certain commands. For example `command=HGET`.
- `percentage`: limits how many times it applies the rule. For example `percentage=25` applies only to 25% of matching commands.
- `client_addr`: scopes the radius to clients coming from a certain subnet,ip,ip:port. For example `client_addr=192.0.0.1`

## Project Status
RedFI is still in its early stages, so expect rapid changes.  
Also, for now, we advise you to use the proxy in development and staging area only.  
If you want to use it in production, you should limit the percentage of connections you send to the proxy.