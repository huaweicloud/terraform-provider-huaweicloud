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

* `region_name` - The region where resources will be created
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the APIG instance
* `plugin_name` - The name of the proxy cache plugin

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zones` - The availability zones to which the instance belongs (default: [])  
  If not specified, will be automatically allocated based on the number of availability_zones_count
* `availability_zones_count` - The number of availability zones to which the instance belongs (default: 1)
* `instance_edition` - The edition of the APIG instance (default: "BASIC")
* `enterprise_project_id` - The ID of the enterprise project, required for enterprise users (default: null)
* `plugin_description` - The description of the proxy cache plugin (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_apig_instance_name"
  plugin_name         = "your_plugin_name"
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

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0 |
| huaweicloud | >= 1.49.0 |
