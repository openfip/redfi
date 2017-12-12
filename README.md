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
RedFI is a proxy that sets between the client and the actual Redis server.
On every incoming command from the client,
it checks the list of failure rules provided by you, and then it applies the first rule that matches the request.

## Usage
First, download the [latest binary](https://github.com/redfi/redfi/releases/tag/v0.1).
Once you're done open the terminal and create a plan.json
```bash
touch ~/plan.json
```

Then copy the following snippet into that file
```javascript
{
    // the failure proxy will pick the first rule that applies to the client request
    // so try to have the specific rules at the top, and the general ones at the bottom
    "rules": [
        // will always delay GET commands by 1000 milliseconds
        {
            "command": "GET",
            "delay": 1000
        },
        // drop connections of client 159.0.147.225:9003 only 50% of the time
        {
            // clientAddr is prefix matched
            // so you could match on subnet level too.
            "clientAddr": "159.0.147.225:9003",
            "percentage": 50,
            "drop": true
        },
        // this rule would be selected for all commands that aren't by 159.0.147.225:9003
        // and aren't GET commands.
        // But it will return nil to the client instead of the real response
        // only 50% of the time
        {
            "percentage": 50,
            "returnEmpty": true
        }
    ]
}
```

After that, execute the binary we downloaded earlier
```bash
$ ./redfi -addr 127.0.0.1:8083 -redis 127.0.0.1:6379 -plan ~/plan.json

RedFI is listening on  127.0.0.1:8083
Don't forget to point your client to that address.
```
- **addr**: is the address on which the proxy listens on for new connections.
- **redis**: address of the actual Redis server to proxy commands/connections to.
- **plan**: path to the json file that contains the rules/scenarios for fault injection.


## Project Status
RedFI is still in its early stages, so expect rapid changes.
Also, for now, we advise you to use the proxy in development and staging area only.
If you want to use it in production, you should limit the percentage of connections you send to the proxy.