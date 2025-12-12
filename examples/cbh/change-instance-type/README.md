# Modify the instance type of the single-node CBH instance

This example provides best practice code for using Terraform to modify the instance type of the single-node CBH instance
in HuaweiCloud.

## Prerequisites

* A Huawei Cloud account
* Terraform installed
* Huawei Cloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the CBH instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `server_id` - The ID of the single node CBH instance to change type

#### Optional Variables

* `availability_zone` - The availability zone of the single-node CBH instance to change type (default: "")

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  server_id = "your_server_id"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* It takes about 15 minutes to modify the instance type of the single-node CBH instance
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.80.4 |
