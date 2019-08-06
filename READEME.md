# Heroku HTTP Adapter

This adapter proxies traffic to the command while unwrapping the cloudevent v0.3 envelope.

## Usage

`PORT=9000 heroku-http-adapter node server.js`

A request that looks like this:

```
curl -X POST \
     -d'@../payload/data-0.json' \
     -H'Content-Type:application/cloudevents' \
     -H'ce-datacontenttype:application/json' \
     -H'ce-specversion:0.3' \
     -H'ce-type:com.github.pull.create' \
     -H'ce-source:https://github.com/cloudevents/spec/pull/123' \
     -H'ce-id:45c83279-c8a1-4db6-a703-b3768db93887' \
     -H'ce-time:2019-06-21T17:31:00Z' \
     http://this-proxy:9000/
```

Will become:

```
curl -X POST \
     -d'@../payload/data-0.json' \
     -H'Content-Type:application/json' \
     -H'ce-specversion:0.3' \
     -H'ce-type:com.github.pull.create' \
     -H'ce-source:https://github.com/cloudevents/spec/pull/123' \
     -H'ce-id:45c83279-c8a1-4db6-a703-b3768db93887' \
     -H'ce-time:2019-06-21T17:31:00Z' \
     http://this-proxys-target:8080/
```
