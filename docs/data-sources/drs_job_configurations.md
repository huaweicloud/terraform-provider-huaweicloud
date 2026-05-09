---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_configurations"
description: |-
  Use this data source to get the list of job parameter configurations for a specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_job_configurations

Use this data source to get the list of job parameter configurations for a specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_job_configurations" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `name` - (Optional, String) Specifies the parameter name to filter the results.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `parameter_config_list` - The list of DRS job parameter configurations.

  The [parameter_config_list](#parameter_config_list_struct) structure is documented below.

<a name="parameter_config_list_struct"></a>
The `parameter_config_list` block supports:

* `name` - The name of the parameter.

* `value` - The current value of the parameter.

* `default_value` - The default value of the parameter.

* `value_range` - Indicates the value range.

* `is_need_restart` - Whether the job needs to be restarted for the parameter modification to take effect.

* `description` - The description of the parameter.

* `created_at` - The creation time of the parameter, in RFC3339 format.

* `updated_at` - The last update time of the parameter, in RFC3339 format.
