---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_parameters_compare"
description: |-
  Manages a TaurusDB HTAP StarRocks parameters compare resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_parameters_compare

Manages a TaurusDB HTAP StarRocks parameters compare resource within HuaweiCloud.

This is a one-time action resource that compares parameters between the source parameter template and the target
parameter template of a HTAP StarRocks instance, and returns the differences.  Deleting this resource will not change
the current configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "source_configuration_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_parameters_compare" "test" {
  source_configuration_id = var.source_configuration_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to compare the HTAP StarRocks parameters.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_configuration_id` - (Required, String, NoneUpdatable) Specifies the ID of the source parameter template
  for comparison.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `differences` - The list of parameter differences between the source and target parameter templates.
  The [differences](#htap_starrocks_parameters_compare_differences_attr) structure is documented below.

<a name="htap_starrocks_parameters_compare_differences_attr"></a>
The `differences` block supports:

* `parameter_name` - The parameter name.

* `source_value` - The parameter value in the source parameter template.

* `target_value` - The parameter value in the target parameter template.
