# Batch Associate Tags to Resources

This example provides best practice code for using Terraform to batch associate tags to multiple resources in
HuaweiCloud Tag Management Service (TMS).

## Prerequisites

* A HuaweiCloud account with TMS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Resources that need to be tagged (e.g., DCS instances, ECS instances, etc.)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `associated_resources_configuration` - The configuration of the associated resources
  + `type` - resource type
  + `id` - resource ID
* `resource_tags` - The tags to be associated with the resources

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  associated_resources_configuration = [
    {
      type = "dcs"
      id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    },
    ...
  ]

  resource_tags = {
    foo   = "bar"
    owner = "terraform"
  }
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

* To clean up the resources (remove tags):

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The `project_id` is automatically obtained from the region name using the `huaweicloud_identity_projects` data source
* All resources must exist before tagging them
* Resource types must be valid TMS-supported types.
* The `tags` parameters of this resource and each service resource will affect each other
* It is recommended to manage tags in only one way to avoid conflicts
* You can use `lifecycle.ignore_changes` in service resources to ignore tag changes in corresponding resources:

  ```hcl
  resource "huaweicloud_dcs_instance" "example" {
    # ... other configuration ...

    lifecycle {
      ignore_changes = [
        tags
      ]
    }
  }
  ```

* When deleting this resource, all tags managed by it will be removed from the associated resources
* Tag keys are case-sensitive
* Each resource can have up to `10` tags (include configured tags in the corresponding resource definition)
* The same tag key can have different values on different resources

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.57.0 |
