# Deploy a SecMaster playbook

This example provides best practice code for using Terraform to deploy a SecMaster playbook in HuaweiCloud
for orchestrating security incident response workflows.
This example demonstrates the complete SecMaster playbook deployment workflow with the following steps:

1. **Create Playbook** - Creates the base playbook structure
2. **Create Version** - Creates a version of the playbook for management
3. **Create Rule** - Defines the trigger conditions for the playbook
4. **Create Playbook Action** - Defines the actions to be executed when triggered
5. **Approve Version** - Approves the playbook version for deployment
6. **Review Playbook** - Reviews the complete playbook configuration
7. **Enable Playbook** - Activates the playbook for production use

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* SecMaster workspace already exists and is configured

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SecMaster playbook is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `playbook_name` - The name of the SecMaster playbook
* `rule_conditions` - The condition rule list of the playbook (must contain at least 2 conditions)

#### Optional Variables

* `workspace_id` - The ID of the SecMaster workspace (at least one of workspace_id and workspace_name must be provided)
* `workspace_name` - The name of the SecMaster workspace (at least one of workspace_id and workspace_name must be provided)
* `rule_expression_type` - The expression type of the playbook rule (default: "common")
* `approval_content` - The approval content for the playbook version (default: "Approved for production use")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  workspace_name       = "your_workspace_name"
  playbook_name        = "your_playbook_name"
  rule_expression_type = "custom"
  rule_conditions      = [
    {
      name = "condition1",
      detail = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      data = [
        "environment.domain_id",
        "==",
        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      ]
    },
    {
      name = "condition2",
      detail = "cn-xxx-x",
      data = [
        "environment.region_id",
        "==",
        "cn-xxx-x",
      ]
    }
  ]
  approval_content     = "Approved for production use"
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
* This example creates a complete SecMaster playbook workflow including playbook, version, rule, action, and approval
* The SecMaster workspace must already exist and be properly configured
* At least one of workspace_id and workspace_name must be provided
* The rule_conditions must contain at least 2 conditions for proper playbook configuration
* The playbook resources can be updated only after the playbook is disabled
* The playbook version and playbook action resources can be updated only after the version is deactivated
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.72.1 |
