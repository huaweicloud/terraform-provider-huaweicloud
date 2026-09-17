---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_check_filter_rules"
description: |-
  Manages a resource to check the DRS data filter rules within HuaweiCloud.
---

# huaweicloud_drs_check_filter_rules

Manages a resource to check the DRS data filter rules within HuaweiCloud.

-> This resource is a one-time action resource used to check the DRS data filter rules. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_check_filter_rules" "test" {
  job_id           = var.job_id
  data_filter_type = "common"
  source           = "job"

  general_filtering_list {
    filter_object_list {
      id          = "test_drs-*-*-mysql1"
      parent_id   = "test_drs"
      object_name = "mysql1"
    }

    source_filter_condition = "id > 1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job for which the data filter rules are to be checked.

* `data_filter_type` - (Optional, String) Specifies the data filter type.  
  The valid values are as follows:
  + **common**: Simple condition filtering.
  + **config**: Associated table filtering.

* `source` - (Optional, String) Specifies the source scenario of the validation request.  
  The valid values are as follows:
  + **job**: The filter rule validation during the data synchronization task execution.
  + **compare**: The filter rule validation during the data comparison task execution.

* `general_filtering_list` - (Optional, List) Specifies the list of the general filtering rules.
  This parameter is required when the `data_filter_type` is set to **common**.  
  The [general_filtering_list](#drs_check_filter_rules_general_filtering_list) structure is documented below.

* `advanced_setting_list` - (Optional, List) Specifies the list of the advanced setting rules.
  This parameter is required when the `data_filter_type` is set to **config**.  
  The [advanced_setting_list](#drs_check_filter_rules_advanced_setting_list) structure is documented below.

<a name="drs_check_filter_rules_general_filtering_list"></a>
The `general_filtering_list` block supports:

* `filter_object_list` - (Optional, List) Specifies the list of the filtering objects.
  The [filter_object_list](#drs_check_filter_rules_filter_object_list) structure is documented below.

* `source_filter_condition` - (Optional, String) Specifies the filter condition SQL statement of the source database.
  Only the content after **WHERE** (without **WHERE** and semicolon) needs to be entered, supporting up to **512**
  characters, for example, **id > 1**.

<a name="drs_check_filter_rules_filter_object_list"></a>
The `filter_object_list` block supports:

* `id` - (Required, String) Specifies the unique identifier of the filtering object, in hierarchical path format.
  The format can be one of the following: **库名-\*-*-\*-表名**, **schema名-\*-*-\*-表名** or
  **库名-\*-*-\*-schema名-\*-*-\*-表名**.

* `parent_id` - (Required, String) Specifies the parent object identifier of the filtering object, which can be a
  database name, schema name or **库名-\*-*-\*-schema名**.

* `object_name` - (Required, String) Specifies the name of the filtering object, that is, the table name to be filtered.

* `object_alias_name` - (Optional, String) Specifies the mapped name of the filtering object.

* `object_type` - (Optional, String) Specifies the type of the filtering object, such as table, view, stored procedure.

<a name="drs_check_filter_rules_advanced_setting_list"></a>
The `advanced_setting_list` block supports:

* `db_name` - (Required, String) Specifies the database name.

* `table_name` - (Required, String) Specifies the table name.

* `col_names` - (Required, String) Specifies the column name.

* `prim_key_or_indexes` - (Required, String) Specifies the primary key or unique index, for example, **id** or
  **id,name**.

* `indexes` - (Required, String) Specifies the index used to optimize the query, for example, **id** or **id,name**.

* `values` - (Required, String) Specifies the filter condition, for example, **id > 1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `query_id` - The query ID of the asynchronous check task, which can be used to query the check result.
