---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_permissions"
description: |-
  Use this data source to get the list of the CC authorized instances.
---

# huaweicloud_cc_permissions

Use this data source to get the list of the CC authorized instances.

## Example Usage

```hcl
variable "permission_id" {}

data "huaweicloud_cc_permissions" "test" {
  permission_id = var.permission_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `permission_id` - (Optional, String) Specifies the permission ID.

* `name` - (Optional, String) Specifies the permission name.

* `description` - (Optional, String) Specifies the permission description.

* `cloud_connection_id` - (Optional, String) Specifies the cloud connection ID.

* `instance_id` - (Optional, String) Specifies the network instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permissions` - The list of the authorized instances.

  The [permissions](#permissions_struct) structure is documented below.

<a name="permissions_struct"></a>
The `permissions` block supports:

* `id` - The authorized instance ID.

* `name` - The authorized instance name.

* `description` - The authorized instance description.

* `domain_id` - The account ID to which the authorized instance belongs.

* `project_id` - The project ID to which the authorized instance belongs.

* `region_id` - The region ID to which the authorized instance belongs.

* `status` - The authorized instance status.

* `created_at` - The creation time.

* `instance_id` - The network instance ID that another account allows you to use.

* `instance_type` - The type of the network instance that another account allows you to use.

* `instance_domain_id` - The account ID of the network instance that another account allows you to use.

* `cloud_connection_id` - The cloud connection ID.
