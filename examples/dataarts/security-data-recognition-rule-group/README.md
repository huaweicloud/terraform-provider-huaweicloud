# Create a DataArts Security data recognition rule group

This example provides best practice code for using Terraform to create DataArts Security data recognition rules and
group them into a data recognition rule group within HuaweiCloud DataArts Studio.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DataArts Studio workspace
* At least one data category under the workspace, or leave `category_ids` empty to auto-select from the workspace

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

Configure authentication variables in `authentication.auto.tfvars` (recommended) or `terraform.tfvars`:

* `region_name` - The region where the DataArts Studio instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Conditionally Required Variables

* `workspace_id` - The ID of the workspace. Required when you already know the target workspace ID
* `instance_id` - The ID of the DataArts Studio instance. Required when `workspace_id` is omitted, used together
  with the `huaweicloud_dataarts_studio_workspaces` data source to resolve the workspace

  -> Configure either `workspace_id` or `instance_id`. When `workspace_id` is provided, the workspace data source is
  skipped and `instance_id` is not used.

* `category_ids` - The ID list of data categories used by data recognition rules. Required when the workspace has no
  available data categories to auto-select. If omitted, the first `data_recognition_rule_count` categories returned by
  `huaweicloud_dataarts_security_data_categories` are used

#### Optional Variables

* `workspace_name` - The name of the workspace used to filter results (default: ""). Only effective when `workspace_id`
  is omitted
* `data_recognition_rule_count` - The number of data recognition rules to create and include in the rule group
  (default: 1)
* `rule_group_name` - The name of the data recognition rule group
* `rule_group_description` - The description of the data recognition rule group (default: "")

## Architecture Overview

This example follows a data recognition rule group creation workflow with conditional data source lookups:

1. **Resolve workspace** (optional):
   + When `workspace_id` is omitted, query `huaweicloud_dataarts_studio_workspaces` by `instance_id` (and optionally
     `workspace_name`) to obtain the target workspace ID
2. **Query data categories**:
   + Query `huaweicloud_dataarts_security_data_categories` under the target workspace to resolve category IDs when
     `category_ids` is omitted
3. **Prepare data secrecy levels and recognition rules**:
   + Create `huaweicloud_dataarts_security_data_secrecy_level` resources for each rule
   + Create `huaweicloud_dataarts_security_data_recognition_rule` resources with type `CUSTOM` and method `NONE`
4. **Create rule group**:
   + Group the created rules using `huaweicloud_dataarts_security_data_recognition_rule_group`
5. **Verify rule group metadata**:
   + Query `huaweicloud_dataarts_security_data_recognition_rule_groups` by group name to verify the created group

## Usage

* Copy this example to your working directory.

* Create an `authentication.auto.tfvars` file for credentials:

  ```hcl
  region_name = "your-region-name"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  ```

* Create a `terraform.tfvars` file and fill in the required variables.

  Example A — provide `workspace_id` and `category_ids` directly:

  ```hcl
  workspace_id    = "your-dataarts-studio-workspace-id"
  rule_group_name = "tf_test_rule_group"
  category_ids    = ["your-data-category-id"]
  ```

  Example B — resolve workspace from instance and auto-select categories:

  ```hcl
  instance_id                 = "your-dataarts-studio-instance-id"
  rule_group_name             = "tf_test_rule_group"
  data_recognition_rule_count = 2
  ```

  Example C — create a rule group with description:

  ```hcl
  workspace_id                = "your-dataarts-studio-workspace-id"
  rule_group_name             = "tf_test_rule_group"
  rule_group_description      = "Created by Terraform"
  category_ids                = ["your-data-category-id-1", "your-data-category-id-2"]
  data_recognition_rule_count = 2
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
* `workspace_id` on the rule group is non-updatable. Changing it will recreate the rule group resource.
* `name` and `rule_ids` on the rule group can be updated in place.
* When `workspace_id` is omitted, `instance_id` must be provided or Terraform will fail the workspace data source
  precondition check.
* When `category_ids` is omitted, ensure the workspace has enough data categories for `data_recognition_rule_count`,
  or creation will fail the data recognition rule precondition check.
* The rule group can be imported using `<workspace_id>/<id>`.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.92.0 |
