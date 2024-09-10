---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_backup"
description: ""
---

# huaweicloud_gaussdb_mysql_backup

Manages a GaussDB MySQL backup resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_backup" "test" {
  instance_id = var.instance_id
  name        = "test_backup_name"
  description = "test description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql instance resource. If omitted,
  the provider-level region will be used. Changing this creates a new instance resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

* `name` - (Required, String, ForceNew) Specifies the name of the backup. It must start with a letter and consist of
  `4` to `64` characters. Only letters (case-sensitive), digits, hyphens (-), and underscores (_) are allowed.

* `description` - (Optional, String, ForceNew) Specifies the description of the backup.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `take_up_time` - Indicates the backup duration in minutes.

* `size` - Indicates the backup size in MB.

* `datastore` - Indicates the database information.

  The [datastore](#backup_datastore_struct) structure is documented below.

<a name="backup_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

## Import

The GaussDB Mysql backup can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_backup.test <id>
```
