<img src="https://raw.githubusercontent.com/redfi/redfi/master/static/redfi@2x.png" width="400px">

RedFI acts as a proxy between the client and Redis with the capability
of injecting faults on the fly, based on the rules given by you.

## Features
- Simple to use. It's just a binary that you execute.
- Transparent for the client.
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

```bash
Usage of ./redfi
  -addr string
    	address for the proxy to listen on (default "127.0.0.1:8083")
  -redis string
    	address to the redis server, to proxy requests to (default "127.0.0.1:6379")
  -plan string
    	path to the plan.json
```


### plan.json example
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