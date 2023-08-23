# Prom Health Checker

Simple Go app that just pings a bunch of sites that you (politely) ask it to

## Application Specifics

### Port

The app listens on `tcp/9090`

### Endpoints

| Endpoint Name | What it does                   |
|---------------|--------------------------------|
| `/healthz`    | Health check endpoint          |
| `/metrics`    | Exposes the prometheus Metrics |


## Prometheus metrics

The application exposes 2 useful metrics per monitored app

| Endpoint                    | 
|-----------------------------|
| `http_latency_milliseconds` |
| `http_up`                   |

### `http_latency_milliseconds`

This endpoint responds with the most recent latency in milliseconds for the app to make an `http get` to the site

### `http_up`

This is a binary operator which tells you if the app is reachable (`1`) or not (`0`)

### Example

```text
# HELP http_latency_milliseconds Latency of HTTP requests to endpoints
# TYPE http_latency_milliseconds gauge
http_latency_milliseconds{endpoint="breadnet"} 40
http_latency_milliseconds{endpoint="google"} 329
# HELP http_up Whether the HTTP endpoint is up (1) or down (0)
# TYPE http_up gauge
http_up{endpoint="breadnet"} 1
http_up{endpoint="google"} 1
```

## How to run

```shell
podman run --rm -v $(pwd)/endpoints.yaml:/app/endpoints.yaml europe-west2-docker.pkg.dev/breadnet-container-store/public/prom-status:latest
```
