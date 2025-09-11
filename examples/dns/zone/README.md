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

* `name` - The name of the zone, note the `.` at the end of the name

#### Optional Variables

* `email` - The email address of the administrator managing the zone
* `zone_type` - The type of zone (default: "public")
* `ttl` - The time to live (TTL) of the zone (default: 300)
* `description` - The description of the zone
* `router` - The list of the router of the zone
  - `router_id` - The ID of the associated VPC
  - `router_region` - The region of the VPC
* `enterprise_project_id` - The enterprise project ID of the zone
* `status` - The status of the zone (default: "ENABLE")
* `proxy_pattern` - The recursive resolution proxy mode for subdomains of the private zone (default: "AUTHORITY")
* `tags` - The key/value pairs to associate with the zone
* `dnssec` - Whether to enable DNSSEC for a public zone (default: "DISABLE")

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
