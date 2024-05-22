# Kepler source in aether

## Overview
This document describes how to use the Kepler source plugin in Aether. The aether-kepler-source plugin fetches Kepler's energy consumption data from Prometheus and sends it to Aether for calculating the carbon footprint. Currently the plugin is designed to collect metrics at the container level, but can be extended to collect metrics at whichever level is desired that Kepler supports.

## Prerequisites

1. A running cluster (currently tested with Kubernetes on EKS,GKE, and locally with kind and k3s)
2. [Prometheus][1] operator with an endpoint for metrics collection
3. [Kepler][2] installed and exporting metrics to Prometheus

## Installation

1. Follow the [installation][3] docs to add aether helm repository.

2. Create a `values.yaml` file with the following:

```yaml
# The environment variables to pass to the aether-kepler-source plugin
# what would typically be found in the `.env` file
env:
- name: INTERVAL
  value: 5m
- name: PROVIDER
  value: "aws"
- name: PROMETHEUS_URL
  value: "http://prometheus-server.monitoring.svc.cluster.local" # by default port is 9090

# Plugin information, in this case the name of the plugin and a link
# to the binary
plugins:
  sources:
    - link: https://github.com/re-cinq/aether-kepler-source/releases/download/[version]/kepler
      name: kepler 
```

3. Install the Aether chart with the created values.yaml file

```bash
helm upgrade --install aether aether/aether -f values.yaml
```
__Note__: This will install the aether deployment in the current namespace,
if you want to install it in a different namespace, you can use the `--namespace` flag.

## Troubleshooting

When the plugin is successfully loaded into aether you will see something like this in the logs:

```bash
2024-05-21T12:01:59.524Z [DEBUG] plugin: starting plugin: path=/plugins/sources/kepler args=["/plugins/sources/kepler"]
2024-05-21T12:01:59.525Z [DEBUG] plugin: plugin started: path=/plugins/sources/kepler pid=16
2024-05-21T12:01:59.525Z [DEBUG] plugin: waiting for RPC address: plugin=/plugins/sources/kepler
2024-05-21T12:01:59.534Z [DEBUG] plugin.kepler: {"time":"2024-05-21T12:01:59.534328527Z","level":"Info","msg":"prometheus address","address":"http://prometheus-k8s.monitoring.svc:9090"}
2024-05-21T12:01:59.535Z [DEBUG] plugin: using plugin: version=1
2024-05-21T12:01:59.535Z [DEBUG] plugin.kepler: plugin address: address=/tmp/plugin881951646 network=unix timestamp=2024-05-21T12:01:59.535Z
2024-05-21T12:01:59.537Z [TRACE] plugin.stdio: waiting for stdio data
{"time":"2024-05-21T12:01:59.539163926Z","level":"INFO","msg":"loaded source plugin","plugin":"kepler"}
```

If the plugin doesn't seem to be loading correctly, be sure that the prometheus URL is correct, and that
the [network policies][4] allow ingress traffic to the Prometheus server from aether.

[1]: https://github.com/prometheus-operator/kube-prometheus
[2]: https://sustainable-computing.io/installation/kepler/
[3]: https://aether.green/docs/tutorials/installation/#helm-installation
[4]: https://github.com/prometheus-operator/kube-prometheus/issues/1780#issuecomment-1168854158
