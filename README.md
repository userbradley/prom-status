# Status Checker using Prometheus

## What is this

This is a little experiment I wrote to poll sites and then report their latency and if they're up or not

## How it works

You provide a yaml file called `endpoints.yaml` to the application with the below format

```yaml
checks:
  - name: breadnet
    url: https://breadnet.co.uk
    frequency: 1s
  - name: google
    url: https://google.com
```

If you do not provide `frequency` to the endpoint name it defaults to 30 seconds. 

If you wish to provide a frequency, ensure that it's in seconds/

## Running the app


---

