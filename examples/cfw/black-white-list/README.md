# Configure blacklist and whitelist filtering for Cloud Firewall

This example provides best practice code for using Terraform to configure blacklist and whitelist filtering in HuaweiCloud
CFW service, including IP-level blacklist rules to block specific IP addresses and whitelist rules to allow specific IP
addresses.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A Cloud Firewall instance already created

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CFW firewall is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `fw_instance_id` - The firewall instance ID

#### Optional Variables

* `blacklist_list_type` - The list type of the blacklist rule (default: 4)
* `blacklist_direction` - The direction of the blacklist rule (default: 0)
* `blacklist_protocol` - The protocol type of the blacklist rule (default: 6)
* `blacklist_port` - The destination port of the blacklist rule (default: "22")
* `blacklist_address_type` - The IP address type of the blacklist rule (default: 0)
* `blacklist_address` - The IP address of the blacklist rule (default: "1.1.1.1")
* `blacklist_description` - The description of the blacklist rule (default: "terraform test")
* `whitelist_list_type` - The list type of the whitelist rule (default: 5)
* `whitelist_direction` - The direction of the whitelist rule (default: 0)
* `whitelist_protocol` - The protocol type of the whitelist rule (default: 6)
* `whitelist_port` - The destination port of the whitelist rule (default: "443")
* `whitelist_address_type` - The IP address type of the whitelist rule (default: 0)
* `whitelist_address` - The IP address of the whitelist rule (default: "2.2.2.2")
* `whitelist_description` - The description of the whitelist rule (default: "terraform test")

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  fw_instance_id = "your_firewall_instance_id"
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
* The `fw_instance_id` can be obtained from the CFW console or the `huaweicloud_cfw_firewalls` data source
* The `object_id` is automatically obtained from the firewall data source
* The blacklist rule blocks traffic from the specified IP address
* The whitelist rule allows traffic from the specified IP address
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- | ---- |
| terraform | >= 1.1.0  |
| huaweicloud | >= 1.47.0 |
