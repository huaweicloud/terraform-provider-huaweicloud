---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_replicate_backup"
description: |-
  Manages a resource to replicate CBR backup within HuaweiCloud.
---

# huaweicloud_cbr_replicate_backup

Manages a resource to replicate CBR backup within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the replicated backup,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "backup_id" {}
variable "destination_project_id" {}
variable "destination_region" {}
variable "destination_vault_id" {}
variable "name" {}
variable "description" {}
variable "enable_acceleration" {}

resource "huaweicloud_cbr_replicate_backup" "example" {
  backup_id = var.backup_id

  replicate {
    destination_project_id = var.destination_project_id
    destination_region     = var.destination_region
    destination_vault_id   = var.destination_vault_id
    name                   = var.name
    description            = var.description
    enable_acceleration    = var.enable_acceleration
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `backup_id` - (Required, String, NonUpdatable) Specifies the ID of the backup to be replicated.

* `replicate` - (Required, List, NonUpdatable) Specifies the replication parameter.
  The [replicate](#replicate_struct) structure is documented below.

<a name="replicate_struct"></a>
The `replicate` block supports:

* `destination_project_id` - (Required, String, NonUpdatable) Specifies the ID of the replication destination project.

* `destination_region` - (Required, String, NonUpdatable) Specifies the replication destination region.

* `destination_vault_id` - (Required, String, NonUpdatable) Specifies the ID of the vault in the replication
  destination region.

* `name` - (Optional, String, NonUpdatable) Specifies the replica name.

* `description` - (Optional, String, NonUpdatable) Specifies the replica description.

* `enable_acceleration` - (Optional, Bool, NonUpdatable) Specifies whether to enable the acceleration function to
  shorten the replication time for cross-region replication. If this parameter is not set, the acceleration function
  is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
