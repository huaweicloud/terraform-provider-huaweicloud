---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_execution_steps"
description: |-
  Use this data source to get the list of COC document execution steps.
---

# huaweicloud_coc_document_execution_steps

Use this data source to get the list of COC document execution steps.

## Example Usage

```hcl
variable "execution_id" {}

data "huaweicloud_coc_document_execution_steps" "test" {
  execution_id = var.execution_id
}
```

## Argument Reference

The following arguments are supported:

* `execution_id` - (Required, String) Specifies the work order ID.

* `execution_step_id_list` - (Optional, List) Specifies the execution step IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of execution steps.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `execution_step_id` - Indicates the work order step ID.

* `action` - Indicates the atomic capability action.

* `start_time` - Indicates the work order step start time.

* `end_time` - Indicates the work order step end time.

* `message` - Indicates the work order step execution information.

* `name` - Indicates the work order step name.

* `status` - Indicates the work order step execution status.

* `inputs` - Indicates the step input parameters.

  The [inputs](#data_inputs_and_outputs_struct) structure is documented below.

* `outputs` - Indicates the step output parameters.

  The [outputs](#data_inputs_and_outputs_struct) structure is documented below.

* `properties` - Indicates the additional attributes for the work order step, stored in map format. For example, to
  display the intranet IP address, it would be **{"fixed_ip": "192.168.1.xx"}.**

<a name="data_inputs_and_outputs_struct"></a>
The `inputs` and `outputs` block supports:

* `key` - Indicates the key of the parameter.

* `value` - Indicates the value of the parameter.
