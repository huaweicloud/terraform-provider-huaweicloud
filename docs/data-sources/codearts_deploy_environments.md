---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_environments"
description: |-
  Use this data source to get the list of CodeArts deploy environments.
---

# huaweicloud_codearts_deploy_environments

Use this data source to get the list of CodeArts deploy environments.

## Example Usage

```hcl
variable "project_id" {}
variable "application_id" {}

data "huaweicloud_codearts_deploy_environments" "test" {
  project_id     = var.project_id
  application_id = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the project ID.

* `application_id` - (Required, String) Specifies the application ID.

* `name` - (Optional, String) Specifies the environment name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `environments` - Indicates the environment lists.

  The [environments](#environments_struct) structure is documented below.

<a name="environments_struct"></a>
The `environments` block supports:

* `id` - Indicates the environment ID.

* `name` - Indicates the environment name.

* `os_type` - Indicates the operating system.

* `created_at` - Indicates the created time.

* `created_by` - Indicates the creator information.

  The [created_by](#environments_created_by_struct) structure is documented below.

* `description` - Indicates the environment description.

* `deploy_type` - Indicates the deployment type.
  The value can be as follows:
  + **0**: host
  + **1**: kubernetes

* `instance_count` - Indicates the number of hosts in the environment.

* `permission` - Indicates the user permission.

  The [permission](#environments_permission_struct) structure is documented below.

<a name="environments_created_by_struct"></a>
The `created_by` block supports:

* `user_id` - Indicates the user ID.

* `user_name` - Indicates the user name.

<a name="environments_permission_struct"></a>
The `permission` block supports:

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_deploy` - Indicates whether the user has the deploy permission.

* `can_edit` - Indicates whether the user has the edit permission.

* `can_manage` - Indicates whether the user has the management permission.

* `can_view` - Indicates whether the user has the view permission.
