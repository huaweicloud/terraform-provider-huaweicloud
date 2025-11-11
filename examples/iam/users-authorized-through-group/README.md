# Authorize users through user group

This example provides best practice code for using Terraform to authorize users through IAM user group in HuaweiCloud
IAM service.

## Prerequisites

* A HuaweiCloud account with IAM permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Introduction

This example demonstrates how to authorize users through user group assignment, which is a recommended approach for
managing user permissions in HuaweiCloud. The workflow includes:

1. **Create or use an existing IAM user group**
2. **Assign roles to the user group** (at project or domain level)
3. **Create IAM users** (with optional auto-generated passwords)
4. **Add users to the user group** to grant them permissions through the group

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IAM service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `users_configuration` - The configuration list of IAM users to be created
  - `name` - The name of the IAM user

#### Optional Variables

* `users_configuration.password` - The password of the IAM user. If not provided, a random password will be generated
  automatically
* `role_id` - The ID of the existing IAM role to be assigned (required if `role_name` and `role_policy` are not
  provided)
* `role_name` - The name of the IAM role.  
  When using existing role, it's required if `role_id` is not provided.  
  When creating custom role, it's required along with `role_policy` and `role_description`
* `role_policy` - The policy JSON string of the custom IAM role. When provided, a new custom role will be created
* `role_description` - The description of the custom IAM role (required when `role_policy` is provided)
* `role_type` - The type of the IAM role (default: "XA"). Only used when creating a custom role
* `group_id` - The ID of the existing IAM group to use (required if `group_name` is not provided)
* `group_name` - The name of the IAM group to be created (required if `group_id` is not provided)
* `group_description` - The description of the IAM group
* `authorized_project_id` - The ID of the project where the role will be assigned (required if `authorized_project_name`
  is not provided and `authorized_domain_id` is not provided)
* `authorized_project_name` - The name of the project where the role will be assigned (required if
  `authorized_project_id` is not provided and `authorized_domain_id` is not provided)
* `authorized_domain_id` - The ID of the domain where the role will be assigned. Required when using domain-level
  authorization. Cannot be used together with project-level authorization parameters

#### Notes

* Either project-level or domain-level authorization should be specified, but not both. When using project-level
  authorization, `authorized_domain_id` is not required. When using domain-level authorization, `authorized_domain_id`
  is required.
* When creating a custom role (`role_policy` is provided), `role_description` must be provided.
* The `role_type` parameter is only used when creating a custom role, and defaults to "XA" if not specified.

## Architecture

This example demonstrates the user group authorization workflow:

1. **Query or create IAM role** - Uses existing role by ID or name, or creates a new custom role with policy
2. **Create or use IAM group** - Creates a new group or uses an existing group
3. **Query project information** - If project-level authorization is used, queries project details
4. **Assign role to group** - Grants the role to the user group at the specified scope (project or domain)
5. **Create IAM users** - Creates one or more IAM users with optional password generation
6. **Add users to group** - Associates users with the group to inherit group permissions

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  role_name               = "system_all_3"   # ECS FullAccess
  group_name              = "tf_test_group"
  authorized_project_name = "cn-north-4"
  users_configuration     = [
    {
      name = "tf_test_user"
      # password is optional, will be auto-generated if not provided
    }
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

### Scenario 1: Create new group and users with project-level authorization

```hcl
role_name               = "system_all_3"   # ECS FullAccess
group_name              = "tf_test_group"
authorized_project_name = "cn-north-4"
users_configuration     = [
  {
    name = "tf_test_user"
  }
]
```

### Scenario 2: Use existing group with domain-level authorization

```hcl
group_id             = "your_existing_group_id"
role_name            = "system_all_3"   # ECS FullAccess
authorized_domain_id = "your_domain_id_to_be_authorized"
users_configuration  = [
  {
    name = "tf_test_user"
  }
]
```

### Scenario 3: Create new group and users with a custom authorization

```hcl
role_name               = "tf_test_role"
role_policy             = <<EOT
{
  "Version": "1.1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "obs:*:*"
      ]
    },
    {
      "Effect": "Deny",
      "Action": [
        "obs:object:DeleteObjectVersion",
        "obs:object:DeleteAccessLabel",
        "obs:bucket:DeleteDirectColdAccessConfiguration",
        "obs:object:AbortMultipartUpload",
        "obs:bucket:DeleteBucketWebsite",
        "obs:object:DeleteObject",
        "obs:bucket:DeleteBucketPolicy",
        "obs:bucket:DeleteBucketCustomDomainConfiguration",
        "obs:object:RestoreObject",
        "obs:bucket:DeleteBucket",
        "obs:object:ModifyObjectMetaData",
        "obs:bucket:DeleteBucketInventoryConfiguration",
        "obs:bucket:DeleteReplicationConfiguration",
        "obs:bucket:DeleteBucketTagging"
      ]
    }
  ]
}
EOT
role_description        = "Created by Terraform"
group_name              = "tf_test_group"
authorized_project_name = "cn-north-4"
users_configuration     = [
  {
    name = "tf_test_user"
  }
]
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* When `password` is not provided in `users_configuration`, a random password will be automatically generated using the
  `random_password` resource
* The role assignment can be done at either project level or domain level, but not both
* Users added to the group will inherit all permissions assigned to that group
* If using an existing group (`group_id`), make sure the group exists before running this example
* When creating a custom role, ensure the `role_policy` is a valid JSON string following HuaweiCloud IAM policy format
* The `role_type` parameter accepts values like "XA" (custom role) and other supported role types

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.37.0 |
| random | >= 3.0.0 |
