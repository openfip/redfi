<img src="https://raw.githubusercontent.com/redfi/redfi/master/static/redfi@2x.png" width="400px">

RedFI acts as a proxy between the client and redis with the capability
of injecting faults on the fly, based off the rules given by you.

```bash
./redfi -redis 127.0.0.1:6379 -addr 127.0.0.1:8083 -plan plan.json
```


## Why?
To gain trust in systems that utilize redis as a backing store.


### plan.json example
```json
{
    "rules": [
        {
            "command": "get",
            "delay": 1000
        },
        {
            "percentage": 50,
            "clientAddr": "127.0.0.1",
            "drop": true
        }
    ]
}
```