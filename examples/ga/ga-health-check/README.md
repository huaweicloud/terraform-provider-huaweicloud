# Create a GA Health Check

This example provides best practice code for using Terraform to create a GA (Global Accelerator) health check
within HuaweiCloud. The health check monitors the health of endpoints in an endpoint group.

The example creates the full dependency chain:

* GA Accelerator (with IPV4 and IPV6 IP sets)
* GA Listener (with port ranges)
* GA Endpoint Group (in the backend region)
* GA Health Check (for the endpoint group)

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where GA resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `accelerator_name` - The name of the GA accelerator
* `listener_name` - The name of the GA listener
* `endpoint_group_name` - The name of the GA endpoint group

#### Optional Variables

* `accelerator_description` - The description of the GA accelerator (default: "GA accelerator for health check")
* `ip_area` - The area of the IP address, valid values: CM, CT, CU, EU, AP, AF, ME, GE (default: "CM")
* `listener_protocol` - The protocol of the GA listener, valid values: TCP, UDP (default: "TCP")
* `listener_description` - The description of the GA listener (default: "GA listener for health check")
* `port_from` - The start port of the listener port range (default: 4000)
* `port_to` - The end port of the listener port range (default: 4200)
* `endpoint_group_description` - The description of the GA endpoint group (default: "GA endpoint group for health check")
* `backend_region` - The region where the backend resources are located (default: "cn-south-1")
* `tags` - The tags of the GA accelerator and listener (default: {})
* `health_check_enabled` - Whether to enable health check (default: true)
* `health_check_interval` - The health check interval in seconds, range: 1-60 (default: 10)
* `health_check_max_retries` - The maximum number of retries, range: 1-10 (default: 5)
* `health_check_port` - The port used for health check, range: 1-65535 (default: 8001)
* `health_check_timeout` - The timeout duration of health check in seconds, range: 1-60 (default: 10)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  accelerator_name    = "ga-accelerator-test"
  listener_name       = "ga-listener-test"
  endpoint_group_name = "ga-endpoint-group-test"
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
* The `protocol` is set to "TCP" as it is the only supported value
* The `enabled` parameter is set to `true` to enable health check
* It is recommended that the `timeout` value is less than the `interval` value
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.72.1 |
