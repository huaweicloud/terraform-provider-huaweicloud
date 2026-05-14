---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_get_params"
description: |-
  Use this data source to get the parameter configuration information for specified DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_batch_get_params

Use this data source to get the parameter configuration information for specified DRS jobs within HuaweiCloud.

## Example Usage

```hcl
variable "job_ids" { 
  type = list(string) 
}

data "huaweicloud_drs_batch_get_params" "test" { 
  job_ids = var.job_ids 
  refresh = 1 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_ids` - (Required, List) Specifies the list of DRS job IDs to query.

* `refresh` - (Optional, String) Specifies whether to refresh the database parameters. The value **1** indicates
  refreshing from the database, and **0** indicates retrieving from the cache. Set this to **1** for the first call.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `params_list` - The response body containing the queried database parameters.

The [params_list](#params_list_struct) structure is documented below.

<a name="params_list_struct"></a>
The `params_list` block supports:

* `job_id` - The ID of the DRS job.

* `params` - The list of detailed parameter information.

The [params](#params_struct) structure is documented below.

<a name="params_struct"></a>
The `params` block supports:

* `group` - The parameter group.
  The valid values are as follows:
  + **common**: General parameters.
  + **performance**: Performance parameters.

* `key` - The parameter name.
  The valid values are as follows:
  + **binlog_cache_size**
  + **binlog_stmt_cache_size**
  + **bulk_insert_buffer_size**
  + **character_set_server**
  + **collation_server**

* `source_value` - The value of the source database parameter.

* `target_value` - The value of the target database parameter.

* `compare_result` - The comparison result of the source and target parameters.
  The valid values are as follows:
  + **true**
  + **false**

* `data_type` - The data type of the parameter.
  The valid values are as follows:
  + **figure**
  + **string**

* `value_range` - The valid value range of the parameter.

* `need_restart` - The indication of whether a restart is required for the parameter change to take effect.
  The valid values are as follows:
  + **true**
  + **false**
