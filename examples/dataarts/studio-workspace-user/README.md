# Create a DataArts Studio workspace user

This example provides best practice code for using Terraform to add an IAM user to a DataArts Studio workspace and
assign roles within HuaweiCloud DataArts Studio Management Center.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) with IAM read privileges
* An existing DataArts Studio workspace
* An existing IAM user to be added to the workspace

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the DataArts Studio instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `role_ids` - The role ID list of the workspace user. Valid values:
  + **r00001** - administrator
  + **r00002** - developer
  + **r00003** - operator
  + **r00004** - viewer

Configure either `user_id` or `user_name`:

* `user_name` - The name of the IAM user to be added to the workspace. Required when `user_id` is omitted
* `user_id` - The ID of the IAM user to be added to the workspace. Required when `user_name` is omitted

#### Conditionally Required Variables

* `workspace_id` - The ID of the workspace. Required when you already know the target workspace ID
* `instance_id` - The ID of the DataArts Studio instance. Required when `workspace_id` is omitted, used together
  with the `huaweicloud_dataarts_studio_workspaces` data source to resolve the workspace

  -> Configure either `workspace_id` or `instance_id`. When `workspace_id` is provided, the workspace data source is
  skipped and `instance_id` is not used.

#### Optional Variables

* `workspace_name` - The name of the workspace used to filter results (default: ""). Only effective when `workspace_id`
  is omitted

## Architecture Overview

This example follows a workspace user creation workflow with conditional data source lookups:

1. **Resolve workspace** (optional):
   + When `workspace_id` is omitted, query `huaweicloud_dataarts_studio_workspaces` by `instance_id` (and optionally
     `workspace_name`) to obtain the target workspace ID
2. **Query available roles**:
   + Query `huaweicloud_dataarts_studio_workspace_user_roles` under the target workspace to list assignable roles and
     validate the configured `role_ids` via a resource precondition
3. **Resolve IAM user** (optional):
   + When `user_id` is omitted, query `huaweicloud_identity_users` by `user_name` to obtain the IAM user ID
4. **Create workspace user**:
   + Add the IAM user to the workspace and assign roles using `huaweicloud_dataarts_studio_workspace_user`

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables.

  Example A — provide `workspace_id` and `user_name` directly:

  ```hcl
  workspace_id = "your-dataarts-studio-workspace-id"
  user_name    = "your-iam-user-name"
  role_ids     = ["r00001"]
  ```

  Example B — resolve workspace from instance:

  ```hcl
  instance_id = "your-dataarts-studio-instance-id"
  user_name   = "your-iam-user-name"
  role_ids    = ["r00002"]
  ```

  Example C — provide `user_id` directly:

  ```hcl
  workspace_id = "your-dataarts-studio-workspace-id"
  user_id      = "your-iam-user-id"
  role_ids     = ["r00004"]
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

* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.
* The AK/SK used by the provider must have IAM read privileges to query `huaweicloud_identity_users`.
* `workspace_id` and `user_id` are non-updatable. Changing them will recreate the workspace user resource.
* Role assignments can be updated in place by modifying the `role_ids` variable.
* The configured `role_ids` are validated against the roles returned by
  `huaweicloud_dataarts_studio_workspace_user_roles` before creation.
* When `workspace_id` is omitted, `instance_id` must be provided or Terraform will fail the workspace data source
  precondition check.
* When `user_id` is omitted, `user_name` must be provided or Terraform will fail the IAM user data source
  precondition check.
* The resource can be imported using `<workspace_id>/<user_id>`.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.91.0 |
