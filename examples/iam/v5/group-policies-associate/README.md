# Associate IAM v5 Policies with User Group

This example provides best practice code for using Terraform to associate one or more IAM v5 policies with a user group
in HuaweiCloud IAM service.

## Prerequisites

* A HuaweiCloud account with IAM permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IAM service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `policy_names` - The name list of policies to be associated with the user group
* `group_name` - The name of the user group

#### Optional Variables

* `policy_type` - The type of the policy (default: `"system"`)
* `group_description` - The description of the user group (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  group_name   = "tf_test_group"
  policy_names = [
    "ModelArtsFullAccessPolicy",
    "SCMReadOnlyPolicy"
  ]
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

## Example Scenarios

### Scenario 1: Associate Multiple System Policies with New Group

This is the default scenario. The example will:

1. Query system policies by type
2. Filter policies by the provided policy names
3. Create a new user group
4. Associate all matched policies with the user group

Set the following variables:

```hcl
policy_type       = "system"
policy_names      = [
  "ModelArtsFullAccessPolicy",
  "SCMReadOnlyPolicy"
]
group_name        = "tf_test_group"
group_description = "Test group with multiple policies"
```

### Scenario 2: Associate Single System Policy with New Group

If you want to associate only one system policy with a user group:

1. Query system policies by type
2. Filter to find the specified policy
3. Create a new user group
4. Associate the policy with the user group

Set the following variables:

```hcl
policy_type       = "system"
policy_names      = ["ModelArtsFullAccessPolicy"]
group_name        = "tf_test_group"
group_description = "Test group for ModelArts access"
```

### Scenario 3: Associate Custom Policies with New Group

If you want to associate custom policies with a user group:

1. Query custom policies by type
2. Filter policies by the provided policy names
3. Create a new user group
4. Associate all matched custom policies with the user group

Set the following variables:

```hcl
policy_type       = "custom"
policy_names      = [
  "your_custom_policy_name_1",
  "your_custom_policy_name_2"
]
group_name        = "tf_test_group"
group_description = "Test group with custom policies"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The policies must exist before associating them with a user group
* The example queries policies by type, then filters by the provided policy names in the `policy_names` list
* If a policy in the list is not found, that specific policy association will be skipped
* The `policy_type` can be `"system"` for system policies or `"custom"` for custom policies
* You can associate multiple policies with a single user group by providing multiple policy names in the `policy_names`
  list
* The user group will be created if it doesn't exist
* Once policies are associated with a group, all users in that group will inherit the permissions defined in all
  associated policies
* The policy and group association is non-updatable - you need to destroy and recreate to change the association
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.83.0 |
