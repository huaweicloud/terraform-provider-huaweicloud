---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_checkpoint_copy"
description: |-
  Using this resource to copy a CBR checkpoint within HuaweiCloud.
---

# huaweicloud_cbr_checkpoint_copy

Using this resource to copy a CBR checkpoint within HuaweiCloud.

-> This resource is only a one-time action resource to copy a CBR checkpoint. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.
This resource has usage restrictions, please refer to
[Replicating a Vault Across Regions](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0009.html).

## Example Usage

```hcl
variable "vault_id" {}
variable "destination_project_id" {}
variable "destination_region" {}
variable "destination_vault_id" {}
variable "auto_trigger" {}

resource "huaweicloud_cbr_checkpoint_copy" "test" {
  vault_id               = var.vault_id
  destination_project_id = var.destination_project_id
  destination_region     = var.destination_region
  destination_vault_id   = var.destination_vault_id
  enable_acceleration    = false
  auto_trigger           = var.auto_trigger
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `vault_id` - (Required, String, NonUpdatable) Specifies the ID of the source vault where the backup to be copied is
  located.

* `destination_project_id` - (Required, String, NonUpdatable) Specifies the ID of the destination project to which the
  backup is to be copied.

* `destination_region` - (Required, String, NonUpdatable) Specifies the ID of the destination region to which the backup
  is to be copied.

* `destination_vault_id` - (Required, String, NonUpdatable) Specifies the ID of the destination vault to which the backup
  is to be copied. The protection type of this vault is required to be **replication**.

* `auto_trigger` - (Optional, Bool, NonUpdatable) Specifies whether to automatically trigger the replication.

* `enable_acceleration` - (Optional, Bool, NonUpdatable) Specifies whether to enable acceleration for cross-region
  replication.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
