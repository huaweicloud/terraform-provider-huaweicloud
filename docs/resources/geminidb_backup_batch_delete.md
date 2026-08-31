---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_backup_batch_delete"
description: |-
  Manages a resource to batch delete manual backup within HuaweiCloud.
---

# huaweicloud_geminidb_backup_batch_delete

Manages a resource to batch delete manual backup within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "backup_id_list" {
  type = list(string)
}

resource "huaweicloud_geminidb_backup_batch_delete" "test" {
  backup_ids = var.backup_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `backup_ids` - (Required, List, NonUpdatable) Specifies the IDs of the backup to be batch delete.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
