---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_auto_scaling_records"
description: |-
  Use this data source to get the list of historical records of auto-scaling.
---

# huaweicloud_gaussdb_mysql_auto_scaling_records

Use this data source to get the list of historical records of auto-scaling.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_mysql_auto_scaling_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the list of records for auto scaling.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - Indicates the record ID.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `scaling_type` - Indicates the scaling type.

* `original_value` - Indicates the original value.

* `target_value` - Indicates the target value.

* `result` - Indicates the scaling result.

* `create_at` - Indicates the scaling time.
