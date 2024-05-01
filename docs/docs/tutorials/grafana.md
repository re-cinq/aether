---
sidebar_position: 4
title: "Setting up Grafana"
---

# Setting up the Grafana dashboard

### Clone Kube Prometheus

**NOTE** if prometheus is already installed you can skip these steps

Deploy kube-prometheus

```bash
git clone --depth 1 https://github.com/prometheus-operator/kube-prometheus; cd kube-prometheus;
```

### Generate the configmap

This step is optional as you can copy the [dashboard json][1] and import
directly in the Grafana UI

```bash
AETHER_EXPORTER_GRAFANA_DASHBOARD_JSON=`curl -fsSL https://raw.githubusercontent.com/re-cinq/aether/main/grafana/aether-exporter.json | sed '1 ! s/^/         /'` 

mkdir -p grafana-dashboards 

cat - > ./grafana-dashboards/aether-exporter-configmap.yaml << EOF 
apiVersion: v1 
kind: ConfigMap 
metadata: 
    labels: 
        app.kubernetes.io/component: grafana 
        app.kubernetes.io/name: grafana 
        app.kubernetes.io/part-of: kube-prometheus 
        app.kubernetes.io/version: 9.5.3 
    name: grafana-dashboard-aether-exporter 
    namespace: monitoring 
data: 
    aether-exporter.json: |- 
      $AETHER_EXPORTER_GRAFANA_DASHBOARD_JSON 
EOF
```

Update the grafana to volume in the new configmap

```bash
 yq -i e '.items += [load("./grafana-dashboards/aether-exporter-configmap.yaml")]' ./manifests/grafana-dashboardDefinitions.yaml 

 yq -i e '.spec.template.spec.containers.0.volumeMounts += [ {"mountPath": "/grafana-dashboard-definitions/0/aether-exporter", "name": "grafana-dashboard-aether-exporter", "readOnly": false} ]' ./manifests/grafana-deployment.yaml 

 yq -i e '.spec.template.spec.volumes += [ {"configMap": {"name": "grafana-dashboard-aether-exporter"}, "name": "grafana-dashboard-aether-exporter"} ]' ./manifests/grafana-deployment.yaml
```


### Deploy Kube Prometheus

Apply the manifests

```bash
kubectl apply --server-side -f manifests/setup 

until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done 

kubectl apply -f manifests/
```

Lastly you will need to apply a service monitor to get prometheus to scrape
aether

```bash
cat <<EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: aether
    app.kubernetes.io/name: aether
    app.kubernetes.io/part-of: kube-prometheus
  name: aether
spec:
  endpoints:
  - interval: 5m
    path: /metrics
    port: http
  jobLabel: app
  selector:
    matchLabels:
      app.kubernetes.io/name: aether
EOF
```

You should now get the dashboards in grafana

```bash
GRAFANA_POD=$(
    kubectl get pod \
        -n monitoring \
        -l app.kubernetes.io/name=grafana \
        -o jsonpath="{.items[0].metadata.name}"
)

k -n monitoring port-forward $GRAFANA_POD 3000
```

You can navigate to grafana on [localhost:3000][2] and using the default user
(`admin`) and password (`admin`), you can see the aether dashboard 

[1]: https://raw.githubusercontent.com/re-cinq/aether/main/grafana/aether-exporter.json
[2]: http://localhost:3000/d/a6dcb47e-4c8c-4295-8c02-087fb01a25a2/aether-exporter
