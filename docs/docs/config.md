---
sidebar_position: 6
title: "Configuration"
---

# Configuration

Aether supports passing the various configurations via a config file.

## Location
The config file will be loaded from the following locations and in the following order:
1. the working directory
2. a sub-directory: `./conf`

## Naming
The name of the file is by default: `local.yaml`  
This helps with testing out the software locally.  
The name of the config file can be changed by setting the following environment variable:  

`AETHER_CONFIG=aether`  

In the example above the new config file name will be: `aether.yaml`


## Configuration
| value                                            | description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | default            |
|--------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------|
| api.address                                      | The IP address to run on                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           | 0.0.0.0            |
| api.port                                         | The port that aether runs on, please make sure to set the correct values for the containerPort  and service if this value changes in kubernetes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | 8080               |
| api.metricsPath                                  | The path where you want to serve metrics from                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | /metrics           |
| plugins.exporterDir                              | The path to load any exporter binary plugins from                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | /plugins/exporters |
| plugins.sourceDir                                | The path to load any source binary plugins from                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | /plugins/sources   |
| providersConfig.scrapingInterval                 | How often aether should fetch metrics and calculate emissions, note its recommended to keep this value unless you understand the implications                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | 5m                 |
| providers.gcp.accounts                           | Google Cloud provider account configuration                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | []                 |
| providers.gcp.accounts.0.project                 | The google cloud project to scrape metrics for                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | null               |
| providers.gcp.accounts.0.credentials.0.filePaths | The credentials used to scrape the projects,  defaults to look for GOOGLE_APPLICATION_CREDENTIALS                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | null               |
| providers.aws.regions                            | List of regions to read the cloud watch metrics for                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | []                 |
| providers.aws.namespaces                         | A namespace is a container for CloudWatch metrics.  Metrics in different namespaces are isolated from each other,  so that metrics from different applications are not mistakenly aggregated into the same statistics. https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | []                 |
| providers.aws.credentials                        | If the credentials config is empty then, aether will try use the aws sdk default  credentials chain: 1. Environment variables.   a. Static Credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)   b. Web Identity Token (AWS_WEB_IDENTITY_TOKEN_FILE) 2. Shared configuration files.   a. SDK defaults to credentials file under .aws folder that is placed in the home folder       on the computer.   b. SDK defaults to config file under .aws folder that is placed in the home folder       on the computer. 3. If your application uses an ECS task definition or RunTask API operation,     IAM role for tasks. 4. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.  Otherwise you can specify one or more locations where to look for either the credentials  or the config or both | []                 |
| providers.aws.credentials.0.profile              | The specific profile to load the credentials for                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | default            |
| providers.aws.credentials.0.filePaths            | The file paths where the credentials file is located                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                | []                 |
| providers.aws.config                             | Load the config for the specific profile, if not set it uses the [default] profile.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |                    |
| providers.aws.config.profile                     | The profile to use                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | default            |
| providers.aws.config.filePaths                   | The file paths where the profile is located                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | []                 |
## Example

```YAML
# Defines the address and port the aether api listens to
api:
  # The address the API server should listen to
  # Can be overridden via: AETHER_API_ADDRESS=localhost
  # Default: 127.0.0.1
  address: 127.0.0.1

  # The port the API server should listen to
  # Can be overridden via: AETHER_API_PORT=8181
  # Default: 8080
  port: 8080

# Aether can use a proxy if necessary
# IMPORTANT: if set, the proxy configuration is applied to all providers
proxy:
  # Can be overridden via: AETHER_PROXY_HTTP_PROXY=http://localhost:3128
  httpProxy: 'http://squid:3128'

  # Can be overridden via: AETHER_PROXY_HTTPS_PROXY=https://localhost:3128
  httpsProxy: 'https://squid:3128'

  # Can be overridden via: AETHER_PROXY_NO_PROXY=localhost
  noProxy: 'intranet.example.com'

# Aether support pulling data from multiple providers
# Each provider has a set of configurations
providers:
  # Google Cloud Provider
  gcp:
    accounts:
      # List of projects to scrape
      - project: 'my-google-cloud-project-id'
        credentials: # optional, defaults to GOOGLE_APPLICATION_CREDENTIALS
          filePaths:
            - '/credentials/application_default_credentials.json'
  # AWS Provider  
  aws:
    # List of regions to read the cloud watch metrics for
    regions:
      - us-east-2
      - us-west-1

    # A namespace is a container for CloudWatch metrics. 
    # Metrics in different namespaces are isolated from each other, 
    # so that metrics from different applications are not mistakenly aggregated into the same statistics.
    # https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/aws-services-cloudwatch-metrics.html
    namespaces:
      - 'AWS/EC2' # EC2
      - 'ContainerInsights' # EKS

    # If the credentials config is empty then, aether will try use the aws sdk default 
    # credentials chain:
    # 
    # 1. Environment variables.
    #    a. Static Credentials (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)
    #    b. Web Identity Token (AWS_WEB_IDENTITY_TOKEN_FILE)
    # 2. Shared configuration files.
    #    a. SDK defaults to credentials file under .aws folder that is placed in the home folder
    #       on the computer.
    #    b. SDK defaults to config file under .aws folder that is placed in the home folder 
    #       on the computer.
    # 3. If your application uses an ECS task definition or RunTask API operation, 
    #    IAM role for tasks.
    # 4. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

    # Otherwise you can specify one or more locations where to look for either the credentials 
    # or the config or both    
    credentials:
      # Load the credentials for the specific profile, if not set it uses the [default] profile. 
      # Example:
      # [default]
      # aws_access_key_id = <YOUR_ACCESS_KEY_ID>
      # aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>      
      #
      profile: 'profile_name'
      filePaths: 
        - 'full_credentials_file_path'
      
    config:
      # Load the config for the specific profile, if not set it uses the [default] profile.
      # Example:
      # [default]
      # region = <REGION>
      profile: 'profile_name'
      filePaths:
        - 'full_file_path'

    # Allows to configure various TCP parameters for the connection to the AWS API
    transport:
      # This setting represents the maximum amount of time to keep an idle network connection 
      # alive between HTTP requests.
      # Set to 0 for no limit.
      # See https://golang.org/pkg/net/http/#Transport.IdleConnTimeout
      # Valid time units are: "ms", "s", "m"
      #
      # Default is zero
      idleConnTimeout: 5s

      # This setting represents the maximum number of idle (keep-alive) connections across all hosts.
      # One use case for increasing this value is when you are seeing many connections in a 
      # short period from the same clients
      # 0 means no limit.
      # See https://golang.org/pkg/net/http/#Transport.MaxIdleConns
      # Default is zero
      maxIdleConns: 0

      # This setting represents the maximum number of idle (keep-alive) connections
      # to keep per-host.
      # One use case for increasing this value is when you are seeing many connections
      # in a short period from the same clients
      #
      # Default is two idle connections per host.
      #
      # Set to 0 to use DefaultMaxIdleConnsPerHost (2).
      #
      # See https://golang.org/pkg/net/http/#Transport.MaxIdleConnsPerHost
      # Default is zero
      maxIdleConnsPerHost: 0

      # This setting represents the maximum amount of time to wait for a client to
      # read the response header.
      # If the client isn’t able to read the response’s header within this duration,
      # the request fails with a timeout error.
      # Be careful setting this value when using long-running Lambda functions,
      # as the operation does not return any response headers until the Lambda
      # function has finished or timed out.
      # However, you can still use this option with the ** InvokeAsync** API operation.
      #
      # Default is no timeout; wait forever.
      # 
      # See https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout
      # Valid time units are: "ms", "s", "m"
      # Default is zero
      responseHeaderTimeout: 10s

      # This setting represents the maximum amount of time waiting for a
      # TLS handshake to be completed.
      #
      # Zero means no timeout.
      #
      # See https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout
      # Valid time units are: "ms", "s", "m"
      # Default is 10 seconds.
      tlsHandshakeTimeout: 10s
```
