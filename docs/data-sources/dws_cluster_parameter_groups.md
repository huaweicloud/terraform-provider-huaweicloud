---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_parameter_groups"
description: |-
  Use this data source to get the list of the parameter groups under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_parameter_groups

Use this data source to get the list of the parameter groups under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_parameter_groups" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the parameter groups.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The ID of the parameter group.

* `name` - The name of the parameter group.

* `type` - The type of the parameter group.

* `status` - The status of the parameter group.
  + **In-Sync**: Synchronized.
  + **Applying**: In application.
  + **Peading-Reboot**: Take effect after restart.
  + **Sync-Failure**: Application failed.

* `fail_reason` - The reason why the parameter application failed.
