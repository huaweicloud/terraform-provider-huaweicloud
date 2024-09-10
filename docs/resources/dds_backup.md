---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_backup"
description: |-
  Manages a DDS backup resource within HuaweiCloud.
---

# huaweicloud_dds_backup

Manages a DDS backup resource within HuaweiCloud.

## Example Usage

```hcl
variable "dds_instance_id" {}
variable "name" {}

resource "huaweicloud_dds_backup" "test"{
  instance_id = var.dds_instance_id
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the manual backup name.
  The value must be `4` to `64` characters in length and start with a letter (from A to Z or from a to z).
  It is case-sensitive and can contain only letters, digits (from 0 to 9), hyphens (-), and underscores (_).

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the manual backup description.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `instance_name` - Indicates the name of a DDS instance.

* `datastore` - Indicates the database version.
  The [datastore](#datastore_struct) structure is documented below.

* `type` - Indicates the backup type. Valid value:
  + **Manual**: indicates manual full backup.

* `begin_time` - Indicates the start time of the backup. The format is yyyy-mm-dd hh:mm:ss. The value is in UTC format.

* `end_time` - Indicates the end time of the backup. The format is yyyy-mm-dd hh:mm:ss. The value is in UTC format.

* `status` - Indicates the backup status. Valid value:
  + **BUILDING**: Backup in progress
  + **COMPLETED**: Backup completed
  + **FAILED**: Backup failed
  + **DISABLED**: Backup being deleted

* `size` - Indicates the backup size in KB.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the database version. The value can be **4.2**, **4.0**, or **3.4**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The DDS backup can be imported using the `instance_id` and the `id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_dds_backup.test <instance_id>/<id>
```
