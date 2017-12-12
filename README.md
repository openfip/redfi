# RedFI
Test the resiliency of your application against Redis' failures.

RedFI acts as a proxy between the client and redis with the capability
of injecting faults on the fly, based off the rules given by you.

```bash
./redfi -server 127.0.0.1:8003 -plan plan.json
```


## Why?
To gain trust in systems that utilize redis as a backing store.

Plan contains a list of rules.
Rules apply on two scopes: global, command based
Rule can do:
- Delay
- TurnOff
- Corrupt?
Limited failure with:
- Percentage
- Client IP:PORT


### plan json
```json
{
    "rules": [
        {
            "delay": 1e3,
            "clientIP": "127.0.0.1",
            "command": "get khalid",
            "commandPrefix": "get kha",
            "percentage": 20
        },
        {
            // this will apply to all actions
            // that don't contain hset commands
            "percentage": 10,
            "drop": true
        }
    ]
}
```