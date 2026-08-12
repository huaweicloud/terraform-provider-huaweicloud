# Create a SecMaster Dataspace

This example provides best practice code for using Terraform to create a SecMaster dataspace
within HuaweiCloud. The dataspace belongs to a SecMaster workspace.

The example creates the full dependency chain:

* SecMaster Workspace
* SecMaster Dataspace (in the workspace)

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SecMaster resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `workspace_name` - The name of the SecMaster workspace
* `dataspace_name` - The name of the SecMaster dataspace. The name can only contain English letters,
  digits and hyphens (-), and cannot start or end with a hyphen (-), nor can they appear
  consecutively. Valid length: 5-63

#### Optional Variables

* `workspace_description` - The description of the SecMaster workspace (default: "Created by Terraform")
* `dataspace_description` - The description of the SecMaster dataspace (default: "Created by Terraform")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  workspace_name = "secmaster-workspace-test"
  dataspace_name = "dataspace-test"
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
* The `workspace_id` is referenced from the workspace resource
* The `project_name` of the workspace is set to the region name
* The `dataspace_name` is non-updatable after creation
* The dataspace can be imported using `<workspace_id>/<id>` format

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.76.0 |
