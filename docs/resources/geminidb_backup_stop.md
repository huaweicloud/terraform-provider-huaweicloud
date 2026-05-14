---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_backup_stop"
description: |-
  Manages a GeminiDB backup stop resource within HuaweiCloud.
---

# huaweicloud_geminidb_backup_stop

Manages a GeminiDB backup stop resource within HuaweiCloud.

-> This resource is only a one-time action resource for stopping a backup of GeminDB.
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "backup_id" {}

resource "huaweicloud_geminidb_backup_stop" "test" {
  backup_id = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new rds instance resource.

* `backup_id` - (Required, String) Specifies the ID of the backup to stop.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as the backup ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
