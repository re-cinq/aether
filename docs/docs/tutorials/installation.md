---
sidebar_position: 1
title: "Installation"
---

# Installation

## Deploying Aether in Your Cluster

Aether integrates with current monitoring systems and some data sources in
order to **estimate** the emissions of your systems. The idea behind Aether
relies heavily on the idea of [Energy Proportionality](https://en.wikipedia.org/wiki/Energy_proportional_computing). In essence the higher
the utilization of hardware components, the energy consumption grows
Proportionally. Therefore Aether will pull key metrics from various
[Sources][1]
such as CPU Utilization, Memory Usage, Storage (coming soon) and Networking
(coming soon)


Due to the variability of cloud providers, each cloud provider works slightly
differently so please see the [sources][1] for more information on
the different source types

## Helm Installation

We use helm to make it easier to deploy aether, so to get started run:

```bash
# add our helm repo
helm repo add aether https://repo.aether.green

# get the latest release
helm repo update
```

To install you can now run:

```bash
helm upgrade --install aether aether/aether -f values.yaml
```
> Note you will need to supply a `values.yaml` with relevant configuration

The default installation has no sources setup, therefore you will need to
configure either using helm values

| value      | description                                                                  | default |
|------------|------------------------------------------------------------------------------|---------|
| secretName | The name of the secret to volume into aether, generally used for credentials | ""      |
| config     | The yaml config to setup aether with, please see the [config docs][2]        | {}      |
|            |                                                                              |         |

> Note: If you want metrics to be exported to prometheus you will need to deploy an aether [ServiceMonitor][3].
> More information can be found in the [grafana documentation][4].

[1]: ../sources/sources.md
[2]: ../config.md
[3]: https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/design.md#servicemonitor
[4]: grafana.md
```

