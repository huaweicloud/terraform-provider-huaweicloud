---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_queue_associated_users"
description: |-
  Use this data source to query the users associated with the specified workload queue within HuaweiCloud.
---

# huaweicloud_dws_workload_queue_associated_users

Use this data source to query the users associated with the specified workload queue within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "queue_name" {}

data "huaweicloud_dws_workload_queue_associated_users" "test" {
  cluster_id = var.dws_cluster_id
  queue_name = var.queue_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

* `queue_name` - (Required, String) Specifies the workload queue name bound to the users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - All users that associated with the specified workload queue.

  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - The name of the user.

* `occupy_resource_list` - The list of the resources used by the user to run jobs.

  The [occupy_resource_list](#users_occupy_resource_list_struct) structure is documented below.

<a name="users_occupy_resource_list_struct"></a>
The `occupy_resource_list` block supports:

* `resource_name` - The resource name.

* `resource_value` - The resource value.

* `value_unit` - The resource attribute unit.
