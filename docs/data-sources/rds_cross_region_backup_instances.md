---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_cross_region_backup_instances"
description: |-
  Use this data source to get the list of RDS instances for which cross-region backups are created.
---

# huaweicloud_rds_cross_region_backup_instances

Use this data source to get the list of RDS instances for which cross-region backups are created.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_cross_region_backup_instances" "test" { 
    instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance.

* `name` - (Optional, String) Specifies the name of the instance.

* `source_region` - (Optional, String) Specifies the source backup region.

* `source_project_id` - (Optional, String) Specifies the project ID of the source backup region.

* `destination_region` - (Optional, String) Specifies the region where the cross-region backup is located.

* `destination_project_id` - (Optional, String) Specifies the project ID of the target backup region.

* `keep_days` - (Optional, Int) Specifies the number of days to retain cross-region backups.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backup_instances` - Indicates the list of instances for which cross-region backups are created.

  The [backup_instances](#backup_instances_struct) structure is documented below.

<a name="backup_instances_struct"></a>
The `backup_instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the name of the instance.

* `source_region` - Indicates the source backup region.

* `source_project_id` - Indicates the project ID of the source backup region.

* `destination_region` - Indicates the region where the cross-region backup is located.

* `destination_project_id` - Indicates the project ID of the target backup region.

* `datastore` - Indicates the database information.

  The [datastore](#backup_instances_datastore_struct) structure is documented below.

* `keep_days` - Indicates the number of days to retain cross-region backups.

<a name="backup_instances_datastore_struct"></a>
The `datastore` block supports:

* `version` - Indicates the database engine version.

* `type` - Indicates the database engine.
  Its value can be any of the following and is case-insensitive: **MySQL**, **PostgreSQL**, **SQLServer**, **MariaDB**.
