# Prometheus Pusher to [Atlassian Status Page](statuspage.io)

# How to?
- Build
```
git clone git@github.com:cxnam/prometheus-pusher.git
cd prometheus-pusher
go get -u
go build
```
- Docker
```
....
```

- Config

Copy [example](./config-example.yaml) to `config.yaml`

- Queries

Copy query [example](./queries-example.yaml) to `queries.yaml`

The prometheus expr needs to return a single element vector

Docs Prometheus query: [QUERYING PROMETHEUS](https://prometheus.io/docs/prometheus/latest/querying/basics/)