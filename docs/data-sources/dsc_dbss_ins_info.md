---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_dbss_ins_info"
description: |-
  Use this data source to get the DBSS instance information within HuaweiCloud.
---

# huaweicloud_dsc_dbss_ins_info

Use this data source to get the DBSS instance information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_dbss_ins_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `dbss_instance_info_list` - The list of DBSS audit instance information.

  The [dbss_instance_info_list](#dbss_instance_info_list_struct) structure is documented below.

* `dbss_rds_database_list` - The list of RDS database information.

  The [dbss_rds_database_list](#dbss_rds_database_list_struct) structure is documented below.

<a name="dbss_instance_info_list_struct"></a>
The `dbss_instance_info_list` block supports:

* `instance_id` - The DBSS instance ID.

* `instance_name` - The DBSS instance name.

<a name="dbss_rds_database_list_struct"></a>
The `dbss_rds_database_list` block supports:

* `configured` - Whether the RDS database is configured.

* `db_name` - The database name.

* `id` - The database ID.

* `type` - The database type.
