# DCS Redis custom template

This example provides best practice code for using Terraform to create and manage a DCS Redis custom template.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS custom template is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `source_template_id` - The ID of the source template
* `template_name` - The name of the custom template

#### Optional Variables

* `source_type` - The type of the source template: sys or user (default: "sys")
* `template_description` - The description of the custom template (default: "")
* `template_params` - The template params to override as a map (default: {"timeout" = "200"})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  source_template_id = "16"
  template_name      = "your_template_name"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* This resource does not require a DCS instance, VPC, or subnet
* `template_id` and `source_type` are ForceNew, changing them will recreate the template
* `name`, `description`, and `params` can be updated in-place via PUT request
* `source_type` can be `sys` (system template) or `user` (custom template)
* Use data source `huaweicloud_dcs_template_detail` to find available param names and value ranges
* The resource supports import using the template `id`
* After import, `template_id`, `source_type`, and `params` cannot be read back from the API
* Use `lifecycle { ignore_changes = [template_id, source_type, params] }` after importing

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.57.0 |
