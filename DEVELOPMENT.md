# Development Guide

## Local Environment

We have setup docker compose to work with. As Aether makes use of collecting
metrics from various sources, you need to authenticate to these sources, plese
check the [docker compose](../docker-compose.yaml) file to see where the
credentials are volumed in. 

docker compose mounts are set by default set as:
```yaml
    volumes:
      # Volume for Google Cloud Credentials
      - ~/.config/gcloud/application_default_credentials.json:/credentials/application_default_credentials.json
      # volume for AWS credentials
      - ~/.aws/credentials:/credentials/credentials
```

You will also need to setup the config file. There is an example one in the
root directory so you should be able to just:

```bash
cp local.example.yaml local.yaml
```

Have a look at the `local.yaml` and set any settings that you want configured. 

the local config should have a `credentials` line under each provider to point
to the correct location:
```
  aws:
    accounts:
      - regions:
        ...
        credentials:
          filePaths:
            - '/credentials/credentials'
  gcp:
    accounts:
      - projects
        ...
        credentials:
          filePaths:
            - '/credentials/application_default_credentials.json'
```

Once this is setup you should be able to run

```bash
docker compose up
```

and this will start your instance. 

