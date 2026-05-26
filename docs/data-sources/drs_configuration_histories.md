---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_configuration_histories"
description: |-
  Use this data source to get the parameter configuration modification history of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_configuration_histories

Use this data source to get the parameter configuration modification history of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_configuration_histories" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

* `begin_time` - (Optional, String) Specifies the start time in UTC format, for example: 2020-09-01T18:50:20Z.

* `end_time` - (Optional, String) Specifies the end time in UTC format, for example: 2020-09-01T19:50:20Z.

* `name` - (Optional, String) Specifies the parameter name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `parameter_history_config_list` - The parameter configuration modification history list.

  The [parameter_history_config_list](#parameter_history_config_list_struct) structure is documented below.

<a name="parameter_history_config_list_struct"></a>
The `parameter_history_config_list` block supports:

* `name` - The parameter name.

* `old_value` - The old parameter value.

* `new_value` - The new parameter value.

* `is_update_success` - The update result. **true** indicates success, **false** indicates failure.

* `is_applied` - Whether the parameter has been applied. **true** indicates applied, **false** indicates not applied.

* `update_time` - The parameter modification time.

* `apply_time` - The parameter application time.
