---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance_parameters_histories"
description: |-
  Use this data source to get the parameter modification history of a GeminiDB instance.
---

# huaweicloud_geminidb_instance_parameters_histories

Use this data source to get the parameter modification history of a GeminiDB instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_instance_parameters_histories" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the configuration histories.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GeminiDB instance.

* `parameter_name` - (Optional, String) Specifies the parameter name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - The list of parameter modification histories.
  The [histories](#geminidb_instance_parameters_histories) structure is documented below.

<a name="geminidb_instance_parameters_histories"></a>
The `histories` block supports:

* `parameter_name` - The parameter name.

* `old_value` - The old parameter value.

* `new_value` - The new parameter value.

* `update_result` - The update result. Valid values are:
  + **SUCCESS**: Update successful.
  + **FAILED**: Update failed.

* `applied` - Whether the parameter change has been applied.
  + **true**: Applied.
  + **false**: Not applied.

* `updated_at` - The update time in the format **yyyy-MM-ddTHH:mm:ssZ**.

* `applied_at` - The applied time in the format **yyyy-MM-ddTHH:mm:ssZ**.
