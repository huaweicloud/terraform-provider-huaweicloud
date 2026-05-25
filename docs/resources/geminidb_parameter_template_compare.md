---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_parameter_template_compare"
description: |-
  Manages a GeminiDB configuration comparison resource within HuaweiCloud.
---

# huaweicloud_geminidb_parameter_template_compare

Manages a GeminiDB configuration comparison resource within HuaweiCloud.

-> This resource is used to compare the differences between two parameter templates.
This is a one-time operation and cannot be undone.

## Example Usage

```hcl
variable "source_configuration_id" {}
variable "target_configuration_id" {}

resource "huaweicloud_geminidb_parameter_template_compare" "test" {
  source_configuration_id = var.source_configuration_id
  target_configuration_id = var.target_configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `source_configuration_id` - (Required, String, NoneUpdatable) Specifies the ID of the source parameter template to compare.

* `target_configuration_id` - (Required, String, NoneUpdatable) Specifies the ID of the target parameter template to compare.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `differences` - The list of parameter differences between the two configurations.
  The [differences](#geminidb_parameter_template_compare_differences) structure is documented below.

<a name="geminidb_parameter_template_compare_differences"></a>
The `differences` block supports:

* `parameter_name` - The parameter name.

* `source_value` - The parameter value in the source configuration.

* `target_value` - The parameter value in the target configuration.
