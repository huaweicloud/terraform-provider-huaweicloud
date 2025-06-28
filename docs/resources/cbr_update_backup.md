---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_update_backup"
description: |-
  Manages a resource to update CBR backup within HuaweiCloud.
---

# huaweicloud_cbr_update_backup

Manages a resource to update CBR backup within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the updated backup,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "backup_id" {}
variable "name" {}

resource "huaweicloud_cbr_update_backup" "test" {
  backup_id = var.backup_id
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `backup_id` - (Required, String, NonUpdatable) Specifies the backup ID.

* `name` - (Required, String, NonUpdatable) Specifies the backup name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
