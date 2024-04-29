---
sidebar_position: 3
title: "Google Cloud Configuration"
---

# Google Cloud Setup

## Authentication

In order to configure Google Cloud deployment you need to give aether access
to Googles Monitoring. This can be done via two ways

### Provide a Json key file

First you need to be authenticate your [gcloud cli][1] and create a service
account

```bash
gcloud auth login
```

Set the project you want to create the service account for

```bash
# this will use your default configure project
# change to the project you want aether to work with
PROJECT=$(gcloud config get-value project)
```

create the service account
```bash
gcloud iam service-accounts create aether \
  --description "SA used for aether integration" \
  --display-name "aether"
```

aether needs some access to query metrics

```bash
gcloud projects add-iam-policy-binding ${PROJECT} \
  --member serviceAccount:aether@${PROJECT}.iam.gserviceaccount.com \
  --role roles/viewer
```

generate the credentials json file

```bash
gcloud iam service-accounts keys create credentials.json \
    --iam-account=aether@${PROJECT}.iam.gserviceaccount.com
```

and then we can create a kubernetes secret from this file

```bash
kubectl create secret generic credentials \
    -from-file=credentials.json
```

once this secret has been created you can update your `values.yaml`
file pointing to this secret

```yaml
...
secretName: "credentials"
...
```

### Workload Authentication

// TODO


## Provider Configuration

In order to configure aether to pull metrics from a project you need to set the
`gcp` provider configuration as follows in the

```yaml
...
config:
    providers:
    # GCP Provider
    gcp:
        accounts:
        # List of projects to scrape
        - project: 'my-google-cloud-project-id'
          # If credentials is not set, it will default to
          # GOOGLE_APPLICATION_CREDENTIALS 
          credentials:
            # This is where the credentials secret will be 
            # volumed into
            filepath: "/etc/secrets/credentials.json"
...
```

once you have updated your config, you can run:

```bash
helm upgrade --install aether aether/aether -f values.yaml
```

and aether should start scraping your metrics and calculating your emissions

[1]: https://cloud.google.com/sdk/gcloud
