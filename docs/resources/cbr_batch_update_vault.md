---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_batch_update_vault"
description: |-
  Using this resource to batch update CBR vaults within HuaweiCloud.
---

# huaweicloud_cbr_batch_update_vault

Using this resource to batch update CBR vaults within HuaweiCloud.

-> This resource is only a one-time action resource to batch update CBR vaults. Deleting this resource will
not change the current vault configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_cbr_batch_update_vault" "test" {
  smn_notify = true
  threshold  = 80
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to execute the request.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `smn_notify` - (Optional, Bool, NonUpdatable) Specifies whether to enable SMN notification for the vault.

* `threshold` - (Optional, Int, NonUpdatable) Specifies the threshold of the vault capacity in GB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `updated_vaults_id` - The list of vault IDs that have been updated.
