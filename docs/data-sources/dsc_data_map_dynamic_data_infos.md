---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_data_map_dynamic_data_infos"
description: |-
  Use this data source to get the dynamic data infos of data map within HuaweiCloud.
---

# huaweicloud_dsc_data_map_dynamic_data_infos

Use this data source to get the dynamic data infos of data map within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_data_map_dynamic_data_infos" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vpc_id` - (Optional, String) Specifies the VPC ID used to filter resources under a specific VPC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpc_dbss_list` - The database audit instance information list.

  The [vpc_dbss_list](#vpc_dbss_list_struct) structure is documented below.

<a name="vpc_dbss_list_struct"></a>
The `vpc_dbss_list` block supports:

* `dbss` - The dynamic database service list.

  The [dbss](#dbss_struct) structure is documented below.

* `total` - The total number of records.

* `vpc_id` - The VPC ID.

* `vpc_name` - The VPC name.

<a name="dbss_struct"></a>
The `dbss` block supports:

* `dbss_instance_info_list` - The DBSS instance information list.

  The [dbss_instance_info_list](#dbss_instance_info_list_struct) structure is documented below.

* `dbss_rds_database_list` - The DBSS RDS database list.

  The [dbss_rds_database_list](#dbss_rds_database_list_struct) structure is documented below.

<a name="dbss_instance_info_list_struct"></a>
The `dbss_instance_info_list` block supports:

* `instance_id` - The DBSS instance ID.

* `instance_name` - The DBSS instance name.

<a name="dbss_rds_database_list_struct"></a>
The `dbss_rds_database_list` block supports:

* `configured` - Whether the database is configured.

* `db_name` - The database name.

* `id` - The database ID.

* `type` - The database type.
