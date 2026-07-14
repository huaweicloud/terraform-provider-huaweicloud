---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_instances"
description: |-
  Use this data source to query the instance list in the asset catalog within HuaweiCloud.
---

# huaweicloud_dsc_catalog_instances

Use this data source to query the instance list in the asset catalog within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "label_id" {}

data "huaweicloud_dsc_catalog_instances" "test" {
  label_id = var.label_id
}
```

### Filter by instance name

```hcl
variable "label_id" {}
variable "instance_name" {}

data "huaweicloud_dsc_catalog_instances" "test" {
  label_id      = var.label_id
  instance_name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the catalog instances are located.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the ID of the group label to which the instance belongs.

* `type_id` - (Optional, String) Specifies the ID of the data type of the database instance.

-> Exactly one of the `label_id` and `type_id` parameters must be specified.

* `instance_name` - (Optional, String) Specifies the name of the database instance.  
  Fuzzy matching is supported.

* `address` - (Optional, String) Specifies the address of the database instance.  
  Fuzzy matching is supported.

* `user` - (Optional, String) Specifies the access user of the database instance.  
  Fuzzy matching is supported.

* `col_id` - (Optional, String) Specifies the key used to sort the instances.

* `sort` - (Optional, String) Specifies the sorting method for query results.  
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of database instances.  
  The [instances](#catalog_instances_attr) structure is documented below.

<a name="catalog_instances_attr"></a>
The `instances` block supports:

* `instance_id` - The ID of the database instance.

* `instance_name` - The name of the database instance.

* `address` - The address of the database instance.

* `db_infos` - The list of databases that belong to the instance.  
  The [db_infos](#catalog_instances_db_infos) structure is documented below.

* `sensitive_col_num` - The number of sensitive columns in the instance.

* `sensitive_db_num` - The number of sensitive databases in the instance.

* `sensitive_table_num` - The number of sensitive tables in the instance.

* `user` - The access user of the database instance.

<a name="catalog_instances_db_infos"></a>
The `db_infos` block supports:

* `db_id` - The ID of the database.

* `db_name` - The name of the database.

* `db_type` - The type of the database.

* `asset_id` - The ID of the asset to which the database belongs.

* `classifications` - The classification list of the database.

* `latest_scan_time` - The latest scan time of the database, in RFC3339 format.

* `sensitive_level_name` - The name of the sensitive level.

* `color_number` - The color number corresponding to the database sensitive level.

* `sensitive_table_count` - The number of sensitive tables in the database.

* `tags` - The tag list of the database.

* `total_table_count` - The total number of tables in the database.
