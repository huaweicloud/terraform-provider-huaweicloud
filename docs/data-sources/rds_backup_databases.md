---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_backup_databases"
description: |-
  Use this data source to query databases contained in a database level backup in a specified RDS instance.
---

# huaweicloud_rds_backup_databases

Use this data source to query databases contained in a database level backup in a specified RDS instance.

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_id" {}

data "huaweicloud_rds_backup_databases" "test" {
  instance_id = var.instance_id
  backup_id   = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `backup_id` - (Required, String) Specifies the ID of the backup.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bucket_name` - Indicates the name of the backup.

* `database_limit` - Indicates the maximum number of databases that can be restored.

* `databases` - Indicates the list of databases included in the backup.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `database_name` - Indicates the name of the database.

* `backup_file_name` - Indicates the name of the backup file.

* `backup_file_size` - Indicates the size of the backup file, in bytes.
