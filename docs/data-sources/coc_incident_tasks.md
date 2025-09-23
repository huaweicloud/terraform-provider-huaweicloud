---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incident_tasks"
description: |-
  Use this data source to get the list of COC incident tasks.
---

# huaweicloud_coc_incident_tasks

Use this data source to get the list of COC incident tasks.

## Example Usage

```hcl
variable "incident_id" {}

data "huaweicloud_coc_incident_tasks" "test" {
  incident_id = var.incident_id
}
```

## Argument Reference

The following arguments are supported:

* `incident_id` - (Required, String) Specifies the event ticket number.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of incident tasks.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `type` - Indicates the node type.

* `key` - Indicates the node key.

* `name` - Indicates the node name.

* `state` - Indicates the node status.

* `operations` - Indicates the list of operations.

  The [operations](#data_operations_struct) structure is documented below.

<a name="data_operations_struct"></a>
The `operations` block supports:

* `task_id` - Indicates the task ID.

* `key` - Indicates the task key.
