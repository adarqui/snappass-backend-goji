### snappass-backend-goji

Example backend for snappass-core-go using goji.

TODO
----

Still no frontend app. Once I code that, then we just point etc/config.json to those static files.

Usage
-----

```
$ make
$ curl -X POST localhost:5000/pass/poop/day
snap:6adc8e64-c783-41d5-871c-b3d4f0b80748
$ curl localhost:5000/key/snap:6adc8e64-c783-41d5-871c-b3d4f0b80748
poop
$ curl localhost:5000/key/snap:6adc8e64-c783-41d5-871c-b3d4f0b80748
Key not found.
```
