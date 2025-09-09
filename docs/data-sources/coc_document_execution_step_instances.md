---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_execution_step_instances"
description: |-
  Use this data source to get the list of COC document execution step instances.
---

# huaweicloud_coc_document_execution_step_instances

Use this data source to get the list of COC document execution step instances.

## Example Usage

```hcl
variable "execution_step_id" {}

data "huaweicloud_coc_document_execution_step_instances" "test" {
  execution_step_id = var.execution_step_id
}
```

## Argument Reference

The following arguments are supported:

* `execution_step_id` - (Optional, String) Specifies the execution step ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of batch instances.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - Indicates the batch execution instance ID.

* `execution_instance_id` - Indicates the execution instance ID.

* `execution_step_id` - Indicates the execution step ID.

* `start_time` - Indicates the instance execution start time.

* `end_time` - Indicates the instance execution end time.

* `status` - Indicates the instance execution status.

* `message` - Indicates the instance execution information.

* `properties` - Indicates the instance information in a batch.
