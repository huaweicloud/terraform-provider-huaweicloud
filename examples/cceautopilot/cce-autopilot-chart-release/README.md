# Deploy a CCE Autopilot Application

This example provides best practice code for using Terraform to create a CCE autopilot cluster,
upload a Helm chart, and deploy a release in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A Helm chart package (.tgz) file

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Cluster Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `cluster_name` - The CCE autopilot cluster name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `cluster_version` - The Kubernetes version of the cluster (default: "v1.34")
* `cluster_description` - The description of the cluster (default: "")
* `cluster_tags` - The tags of the cluster (default: {})

### Chart Variables

#### Required Variables

* `chart_content` - The path of the chart package to be uploaded

#### Optional Variables

* `chart_parameters` - The parameters of the chart (default: "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}")

### Release Variables

#### Required Variables

* `release_name` - The name of the release
* `release_version` - The version of the release

#### Optional Variables

* `release_namespace` - The namespace to deploy the release (default: "default")
* `release_image_tag` - The image tag (default: "latest")
* `release_image_pull_policy` - The image pull policy (default: "IfNotPresent")
* `release_description` - The description of the release (default: "")
* `release_dry_run` - Whether to dry run (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  cluster_name            = "your_cluster_name"
  chart_content           = "./your-chart-1.0.0.tgz"
  release_name            = "my-release"
  release_version         = "1.0.0"
  release_namespace       = "default"
  release_image_tag       = "v1"
  release_image_pull_policy = "IfNotPresent"
  release_description     = "Deployed by Terraform"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* This example creates a VPC, subnet, CCE autopilot cluster, uploads a Helm chart, and deploys a release
* The creation of the CCE autopilot cluster takes about 10-15 minutes
* CCE autopilot clusters use the `eni` container network mode
* The chart package must be a .tgz file
* The chart content path must be accessible from the Terraform working directory
* The namespace must exist or will be created automatically
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.90.0 |
