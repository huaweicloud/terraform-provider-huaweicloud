---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_configuration_update"
description: |-
  Manages a resource to update DRS job configuration within HuaweiCloud.
---

# huaweicloud_drs_job_configuration_update

Manages a resource to update DRS job configuration within HuaweiCloud.

-> This resource is a one-time action resource used to update DRS job configuration. Deleting this resource
  will not rollback the configuration update operation, but will only remove the resource information from the tf
  state file.

## Example Usage

```hcl
variable "job_id" {}
variable "configuration_params" {
  type = list(object({
    parameter_name  = string
    parameter_value = string
  }))
}

resource "huaweicloud_drs_job_configuration_update" "test" {
  job_id = var.job_id

  dynamic "values" {
    for_each = var.configuration_params

    content {
      parameter_name  = values.value.parameter_name
      parameter_value = values.value.parameter_value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `values` - (Required, List, NonUpdatable) Specifies the list of parameter values to update.
  The [values](#values_struct) structure is documented below.

<a name="values_struct"></a>
The `values` block supports:

* `parameter_name` - (Required, String, NonUpdatable) Specifies the parameter name.
  For example: **applier_thread_num**, **read_task_num**, etc.

* `parameter_value` - (Required, String, NonUpdatable) Specifies the parameter value.
  For example: **20**, **false**, etc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
