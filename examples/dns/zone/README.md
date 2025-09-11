# Create a zone

This example provides best practice code for using Terraform to create a zone in HuaweiCloud DNS service.

## Prerequisites

* A HuaweiCloud account with DNS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DNS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `dns_public_zone_name` - The name of the zone, note the `.` at the end of the name

#### Optional Variables

* `dns_public_zone_email` - The email address of the administrator managing the zone
* `dns_public_zone_type` - The type of zone (default: "public")
* `dns_public_zone_ttl` - The time to live (TTL) of the zone (default: 300)
* `dns_public_zone_description` - The description of the zone
* `dns_public_zone_router` - The list of the router of the zone
  - `router_id` - The ID of the associated VPC
  - `router_region` - The region of the VPC
* `dns_public_zone_enterprise_project_id` - The enterprise project ID of the zone
* `dns_public_zone_status` - The status of the zone (default: "ENABLE")
* `dns_public_zone_dnssec` - Whether to enable DNSSEC for a public zone (default: "DISABLE")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  name = "your_zone_name"
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

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.77.4 |
