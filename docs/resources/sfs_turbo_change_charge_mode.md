---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_change_charge_mode"
description: |-
  Use this resource to change the charge mode of the SFS turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_change_charge_mode

Use this resource to change the charge mode of the SFS turbo within HuaweiCloud.

-> The current resource is a one-time action resource using to change SFS Turbo charge mode. Deleting this resource
will not reset the charge mode, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "share_id" {}
variable "period_num" {}
variable "period_type" {}
variable "is_auto_renew" {}

resource "huaweicloud_sfs_turbo_change_charge_mode" "test" {
  share_id    = var.share_id
  period_num    = var.period_num
  period_type   = var.period_type
  is_auto_renew = var.is_auto_renew
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, NonUpdatable) Specifies the ID of the SFS Turbo.

* `period_num` - (Required, Int, NonUpdatable) Specifies the charging period of the SFS Turbo.

* `period_type` - (Required, Int, NonUpdatable) Specifies the charging period unit of the SFS Turbo.
  + **2**: paid monthly.
  + **3**: paid annual.

* `is_auto_renew` - (Optional, Int, NonUpdatable) Specifies whether auto renew is enabled. Defaults to **0**.
  + **0**: no automatic renewal.
  + **1**: automatic renewal.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
