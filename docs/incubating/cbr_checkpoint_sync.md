---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_checkpoint_sync"
description: |-
  Using this resource to synchronize hybrid cloud checkpoints within HuaweiCloud.
---

# huaweicloud_cbr_checkpoint_sync

Using this resource to synchronize hybrid cloud checkpoints within HuaweiCloud.

-> This resource is only a one-time action resource to synchronize checkpoints from hybrid cloud vaults. Deleting this
resource will not clear the corresponding checkpoint synchronization record, but will only remove the resource
information from the tfstate file.

## Example Usage

```hcl
variable "vault_id" {}
variable "auto_trigger" {}

resource "huaweicloud_cbr_checkpoint_sync" "test" {
  vault_id     = var.vault_id
  auto_trigger = var.auto_trigger
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used.

* `vault_id` - (Required, String, NonUpdatable) Specifies the ID of the hybrid cloud vault to sync checkpoints to.

* `auto_trigger` - (Required, Bool, NonUpdatable) Specifies whether this checkpoint sync is automatically triggered.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
