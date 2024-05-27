---
sidebar_position: 2
title: "AWS Configuration"
---

# AWS Setup

## Authentication

In order to access AWS resources aether will need to authenticate with the
AWS SDK. This is possible with credentials (access key and secret access key).

### Configuration Precedence

To authenticate to the client, aether uses the [default AWS][1] configuration precedence.

1. Environment variables
    a. Static Credentials (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
    b. Web Identity Token (`AWS_WEB_IDENTITY_TOKEN_FILE`)
2. Shared configuration files
    a. SDK defaults to credentials file under .aws folder that is placed in the home folder on your computer.
    b. SDK defaults to config file under .aws folder that is placed in the home folder on your computer.
3. If your application uses an ECS task definition or RunTask API operation, IAM role for tasks.
4. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

Additionally, credentials or a path to a config or credentials file can be added to the local YAML

```yaml
credentials:
  # Load the credentials for the specific profile, if not set it uses the [default] profile.
  # Example:
  # [default]
  # aws_access_key_id = <YOUR_ACCESS_KEY_ID>
  # aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```
> _Note: be careful not to publicly store the config if using this option_

```yaml
credentials:
  profile: 'profile_name'
  filePaths:
    - 'full/credentials/filepath' # Path to credentials file
```
```yaml
config:
  # Load the config for the specific profile, if not set it uses the [default] profile.
  # Example:
  # [default]
  # region = <REGION>
  profile: 'profile_name'
  filePaths:
    - 'full_file_path'
```

## Provider Configuration

To configure aether to pull metrics from `AWS` the YAML needs to be updated with the
provider information.

```yaml
providers:
  # AWS Provider
  aws:
    # List of regions to read the metrics for
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
```

## ServiceAccount Setup

To set up a service account for aether running on Amazon EKS to query the EC2 and CloudWatch metrics API, follow these steps:

The service account needs to have the following policy actions:
* `ec2:DescribeInstances`
* `cloudwatch:GetMetricData`

1. Create an IAM OIDC Identity Provider:
Follow these [steps][4] to see if you have an OIDC identity provider already set up, and if not how to set one up.

2. Create the IAM Policy file:
Create a file named `aether-iam-policy.json` with the following policy:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "cloudwatch:GetMetricData"
      ],
      "Resource": "*"
    }
  ]
}
```

3. Create the IAM Policy:
```bash
aws iam create-policy --policy-name aether-read-metrics --policy-document file://aether-iam-policy.json
```

4. Create and attach the IAM policy to the service account:

By default, the aether helm install creates an empty service account for you, named `aether`

There are multiple ways to attach the IAM policy to the service account:
* If eksctl is installed, you can use the following command to create the IAM role and service account.

```bash
eksctl create iamserviceaccount \
  --name your-service-account \
  --namespace default \
  --cluster your-cluster-name \
  --attach-policy-arn arn:aws:iam::account-id:policy/aether-read-metrics \
  --approve
```
> note: If you created a role for the SA, then attach that with `--role`. Otherwise, eksctl will create a role for you.
  To use the default `aether` service account you need to append the command with `--override-existing-serviceaccounts`

* Follow the [AWS documentation][5] to use the aws client.

* Create a role with the policy via the UI and annotate the service account with the role ARN
```bash
kubectl annotate serviceaccount aether \
  eks.amazonaws.com/role-arn=arn:aws:iam::your-account-id:role/your-role-name
```

5. Validate that the service account has the role annotated. For additional validation tactics, refer to the AWS documentation.
```bash
kubectl get serviceaccount aether -o yaml
```

### CloudWatch

Currently aether only supports metric scraping from [CloudWatch][2], which incurs costs. We are planning to
add OTEL metrics and Kepler as different metric scraping methods in the future. 

Additionally, by default, CloudWatch does *not* collect memory utilization metrics, only those for CPU. 
So to get memory energy consumption, the [CWAgent][3] needs to be installed on instances.


[1]: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk
[2]: https://aws.amazon.com/cloudwatch/
[3]: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Install-CloudWatch-Agent.html
[4]: https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html
[5]: https://docs.aws.amazon.com/eks/latest/userguide/associate-service-account-role.html
