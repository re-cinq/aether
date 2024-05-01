---
sidebar_position: 1
---

# aether

## Introduction

Aether is a calculation engine that uses metrics of infrastructure and
calculates emissions in real-time based on factors. There are three main
components to Aether

- [Sources][1]
- Calculation engine
- Exporters

![Architecture](/img/architecture.webp)

The main design decisions behind Aether are:

- Using metric sources to calculate emissions (we do not install agents, but
  rather pull data from current metric solutions)
- We calculate emissions, not only energy/power. This means that we can use
  sources like [kepler][3] and [scaphandre][2]

[1]: sources/sources.md
[2]: https://github.com/hubblo-org/scaphandre
[3]: https://sustainable-computing.io/
