# Create an APIG instance with proxy cache plugin

This example provides best practice code for using Terraform to create an API Gateway (APIG) instance with a proxy cache
plugin in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `access_key` - HuaweiCloud access key
* `secret_key` - HuaweiCloud secret key
* `region_name` - The region where resources will be created

### Resource Variables

* `vpc_name` - Name of the Virtual Private Cloud (VPC)
* `subnet_name` - Name of the subnet within the VPC
* `security_group_name` - Name of the security group
* `instance_name` - Name of the APIG instance
* `plugin_name` - Name of the proxy cache plugin

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key           = "your_access_key"
  secret_key           = "your_secret_key"
  region_name          = "your_region"
  vpc_name             = "example-vpc"
  subnet_name          = "example-subnet"
  security_group_name  = "example-sg"
  instance_name        = "example-apig"
  plugin_name          = "example-proxy-cache"
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

## Proxy Cache Plugin Configuration

The proxy cache plugin is configured with the following settings:

* **Cache Key Configuration**
  + Custom parameter: `custom_param`
  + No system parameters or headers used in cache key

* **Cache HTTP Status and TTL**
  + HTTP status codes 202 and 203 are cached for 5 seconds

* **Client Cache Control**
  + Mode: off

* **Cacheable Headers**
  + `X-Custom-Header`

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* The APIG instance is created with BASIC edition
* The APIG instance will be deployed in the first available zone
