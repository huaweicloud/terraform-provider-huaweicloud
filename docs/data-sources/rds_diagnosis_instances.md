---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_diagnosis_instances"
description: |-
  Use this data source to obtain the diagnosis result of a specified diagnosis item for RDS instances.
---

# huaweicloud_rds_diagnosis_instances

Use this data source to obtain the diagnosis result of a specified diagnosis item for RDS instances.

## Example Usage

```hcl
variable "engine" {}
variable "diagnosis" {}

data "huaweicloud_rds_diagnosis_instances" "test" {
  engine    = var.engine
  diagnosis = var.diagnosis
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `engine` - (Required, String) Specifies the RDS engine type. The valid values are:
  **mysql**, **postgresql**, **sqlserver**.

* `diagnosis` - (Required, String) Specifies the diagnosis item. The valid values are:
  **high_pressure**, **lock_wait**, **insufficient_capacity**, **slow_sql_frequency**, **age_exceed**,
  **disk_performance_cap**, **mem_overrun**, **connections_exceed**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of instance objects affected by the diagnosis result.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block contains:

* `id` - Indicates the ID of the affected instance.
