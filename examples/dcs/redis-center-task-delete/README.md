# Delete DCS center task

This example provides best practice code for using Terraform to delete a DCS center task in HuaweiCloud DCS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DCS center task that needs to be deleted

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS task is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `task_id` - The ID of the DCS center task to be deleted

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  task_id = "your_dcs_center_task_id"
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
* The `task_id` is non-updatable, changing it will trigger resource replacement
* Deleting this resource only removes it from Terraform state, the task is already deleted by the apply operation

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.94.0 |
