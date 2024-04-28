---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_restore"
description: ""
---

# huaweicloud_dcs_instance_restore

Use this resource to restore a DCS instance with a backup within HuaweiCloud.

-> **NOTE:** Deleting restoration record is not supported. If you destroy a resource of restoration record,
the restoration record is only removed from the state, but it remains in the cloud. And the instance doesn't return to
the state before restoration.

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_id" {}

resource "huaweicloud_dcs_instance_restore" "test" {
  instance_id = var.instance_id
  backup_id   = var.backup_id
  description = "test DCS restoration"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance to be restored.
  Changing this creates a new resource.

* `backup_id` - (Required, String, ForceNew) Specifies the backup ID used to restore the DCS instance.
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the DCS instance restoration.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the restoration record ID.

* `restore_name` - Indicates the name of the restoration record.

* `created_at` - Indicates the time when the restoration record created.

* `updated_at` - Indicates the time when the restoration record completed.

## Timeouts

This resource provides the following timeout configuration option:

* `create` - Default is 30 minutes.
