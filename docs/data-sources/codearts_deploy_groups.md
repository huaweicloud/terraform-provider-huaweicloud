---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_groups"
description: |-
  Use this data source to get the list of CodeArts deploy groups.
---

# huaweicloud_codearts_deploy_groups

Use this data source to get the list of CodeArts deploy groups.

## Example Usage

```hcl
variable "project_id" {}

data "huaweicloud_codearts_deploy_groups" "test" {
  project_id = var.project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the project ID.

* `name` - (Optional, String) Specifies the name of host cluster.

* `os_type` - (Optional, String) Specifies the operating system. Valid values are **windows**, **linux**.

* `is_proxy_mode` - (Optional, String) Specifies whether the host is an agent host.
  Valid values are as follows:
  + **1**: Using proxy access mode.
  + **0**: Without using proxy access mode.

* `resource_pool_id` - (Optional, String) Specifies the customized resource pool ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the host cluster list.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - Indicates the host cluster ID.

* `name` - Indicates the host cluster name.

* `description` - Indicates the description of host cluster.

* `env_count` - Indicates the number of environments.

* `host_count` - Indicates the the number of hosts in a cluster.

* `os_type` - Indicates the operating system.

* `resource_pool_id` - Indicates the slave cluster ID.
  + If the default value is null, the default slave cluster is used.
  + If the value is user-defined, the slave cluster ID is used.

* `permission` - Indicates the permission list.

  The [permission](#groups_permission_struct) structure is documented below.

* `created_by` - Indicates the creator name.

<a name="groups_permission_struct"></a>
The `permission` block supports:

* `can_edit` - Indicates whether the user has the edit permission.

* `can_delete` - Indicates whether the user has the deletion permission.

* `can_add_host` - Indicates whether the user has the permission to add hosts.

* `can_manage` - Indicates whether the user has the management permission.

* `can_view` - Indicates whether the user has the view permission.

* `can_copy` - Indicates whether the user has the permission to copy hosts.
