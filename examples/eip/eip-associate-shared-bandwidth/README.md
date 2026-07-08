# Add Bound EIP to Shared Bandwidth

This example provides best practice code for using Terraform to create a dedicated EIP and add it to shared bandwidth
in HuaweiCloud EIP service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EIP resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bandwidth_name` - The name of the shared bandwidth
* `eip_bandwidth_name` - The name of the dedicated EIP bandwidth

#### Optional Variables

* `enterprise_project_id` - The ID of the enterprise project (default: `""`)
* `bandwidth_size` - The size of the shared bandwidth in Mbit/s (default: 5)
* `bandwidth_charge_mode` - The charge mode of the shared bandwidth (default: "bandwidth")
* `bandwidth_type` - The type of the shared bandwidth (default: "share")
* `bandwidth_public_border_group` - The border group of the public IP (default: "center")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `eip_bandwidth_size` - The size of the dedicated EIP bandwidth in Mbit/s (default: 5)
* `eip_bandwidth_charge_mode` - The charge mode of the dedicated EIP bandwidth (default: "traffic")
* `eip_description` - The description of the EIP (default: `""`)
* `eip_tags` - The tags of the EIP (default: `{}`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bandwidth_name     = "tf_test_shared_bandwidth"
  eip_bandwidth_name = "tf_test_eip_bandwidth"
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

## Shared bandwidth association and disassociation behavior

### Changes after adding the EIP to shared bandwidth

The EIP is added to shared bandwidth through `huaweicloud_eip_bandwidth_associate`, the following changes occur on
the cloud side:

1. `bandwidth.share_type` changes from `PER` to `WHOLE`
2. `bandwidth.id` changes to the shared bandwidth ID
3. The original dedicated bandwidth configuration is no longer used by the EIP

For that scenario, use `lifecycle { ignore_changes = [bandwidth] }` on `huaweicloud_vpc_eip` to avoid perpetual plan
diffs.

### Behavior when destroying the association resource

Destroying only the `huaweicloud_eip_bandwidth_associate` resource will detach the EIP from the shared bandwidth and
restore the EIP to the dedicated bandwidth configuration defined by `bandwidth_charge_mode`, `bandwidth_size`
and `bandwidth_name`.

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The shared bandwidth and EIP must be in the same region
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.92.0 |
