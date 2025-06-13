---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vault_change_charge_mode"
description: |-
  Using this resource to change the charge mode of CBR vaults within HuaweiCloud.
---

# huaweicloud_cbr_vault_change_charge_mode

Using this resource to change the charge mode of CBR vaults within HuaweiCloud.

-> This resource is only a one-time action resource to change the charge mode of CBR vaults. Deleting this resource will
not affect the actual charge mode of the vaults, but will only remove the resource information from the tfstate file.

-> The vaults to change charge mode must be in **postPaid** mode. Using this resource may cause unexpected changes to
the `charging_mode` field of the `huaweicloud_cbr_vault` resource. At this time, you can use the `lifecycle` statement
to ignore the change of `charging_mode`.

## Example Usage

```hcl
variable "vault_id" {}

resource "huaweicloud_cbr_vault_change_charge_mode" "test" {
  vault_ids     = [var.vault_id]
  charging_mode = "pre_paid"
  period_type   = "month"
  period_num    = 1
  is_auto_renew = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `vault_ids` - (Required, List, NonUpdatable) Specifies the IDs of the vaults to change charge mode.

* `charging_mode` - (Required, String, NonUpdatable) Specifies the charging mode of the vault. Only supports **pre_paid**.

* `period_type` - (Required, String, NonUpdatable) Specifies the period type of the vault. Only supports **month** and **year**.

* `period_num` - (Required, Int, NonUpdatable) Specifies the number of periods to purchase.

* `is_auto_renew` - (Optional, Bool, NonUpdatable) Specifies whether to auto-renew the vault when it expires.
  Defaults to **false**.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
