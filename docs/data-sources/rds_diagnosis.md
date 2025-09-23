---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_diagnosis"
description: |-
  Use this data source to query the number of diagnosed RDS instances by diagnosis type for a specific database engine.
---

# huaweicloud_rds_diagnosis

Use this data source to query the number of diagnosed RDS instances by diagnosis type for a specific database engine.

## Example Usage

```hcl
variable "engine" {}

data "huaweicloud_rds_diagnosis" "test" {
  engine = var.engine
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `engine` - (Required, String) Specifies the RDS engine type.
  The valid values are **mysql**, **postgresql**, **sqlserver**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `diagnosis` - Indicates the list of diagnosis result objects.

  The [diagnosis](#diagnosis_struct) structure is documented below.

<a name="diagnosis_struct"></a>
The `diagnosis` block contains:

* `name` - Indicates the diagnosis item.

* `count` - Indicates the number of instances.
