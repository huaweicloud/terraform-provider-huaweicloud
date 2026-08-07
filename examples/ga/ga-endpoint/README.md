# Create a GA Endpoint

This example provides best practice code for using Terraform to create a GA (Global Accelerator) endpoint
within HuaweiCloud. The GA endpoint routes traffic to a backend EIP resource.

The example creates the full dependency chain:

* GA Accelerator (with IPV4 and IPV6 IP sets)
* GA Listener (with port ranges)
* GA Endpoint Group (in the backend region)
* VPC EIP (in the backend region)
* GA Endpoint (pointing to the EIP)

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
* `eip_name` - The name of the EIP bandwidth

#### Optional Variables

* `accelerator_description` - The description of the GA accelerator (default: "")
* `ip_area` - The area of the IP address, valid values: CM, CT, CU, EU, AP, AF, ME, GE (default: "CM")
* `listener_protocol` - The protocol of the GA listener, valid values: TCP, UDP (default: "TCP")
* `listener_description` - The description of the GA listener (default: "GA listener for endpoint")
* `port_from` - The start port of the listener port range (default: 4000)
* `port_to` - The end port of the listener port range (default: 4200)
* `endpoint_group_description` - The description of the GA endpoint group (default: "GA endpoint group")
* `backend_region` - The region where the backend EIP resource is located (default: "cn-south-1")
* `eip_type` - The type of the EIP, valid values: 5_bgp, 5_sbgp (default: "5_bgp")
* `bandwidth_size` - The size of the EIP bandwidth (default: 8)
* `endpoint_weight` - The weight of the endpoint for traffic distribution, range: 0-100 (default: 10)
* `tags` - The tags of the GA accelerator and listener (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  accelerator_name    = "ga-accelerator-test"
  listener_name       = "ga-listener-test"
  endpoint_group_name = "ga-endpoint-group-test"
  eip_name            = "ga-eip-test"
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
* The GA accelerator, listener, endpoint group, and endpoint are created in the `region_name` region
* The EIP is created in the `backend_region` (which may differ from the GA region)
* The `resource_type` is set to "EIP" as it is the only supported value
* The `weight` parameter controls traffic distribution among endpoints in the same group
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.72.1 |
