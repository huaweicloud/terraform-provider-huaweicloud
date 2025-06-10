---
subcategory: "Cloud Backup and Recovery"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vault_migrate_resources"
description: |-
  Use this resource to migrate resources between CBR vaults within HuaweiCloud.
---

# huaweicloud_cbr_vault_migrate_resources

Use this resource to migrate resources between CBR vaults within HuaweiCloud.

-> This resource is only a one-time action resource to migrate resources between CBR vaults. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "source_vault_id" {}
variable "destination_vault_id" {}
variable "resource_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_vault_migrate_resources" "test" {
  vault_id             = var.source_vault_id
  destination_vault_id = var.destination_vault_id
  resource_ids         = var.resource_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the provider-level
  region will be used. Changing this will create a new resource.

* `vault_id` - (Required, String, NonUpdatable) Specifies the source vault ID from which resources will be migrated.

* `destination_vault_id` - (Required, String, NonUpdatable) Specifies the destination vault ID where resources
  will be migrated to.

* `resource_ids` - (Required, List, NonUpdatable) Specifies the IDs of the resources to be migrated.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the source vault (also `vault_id`).

* `migrated_resources` - The list of resources that have been successfully migrated.
