---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_statical_chart"
description: |-
  Use this data source to get the list of DSC catalog statical chart within HuaweiCloud.
---

# huaweicloud_dsc_catalog_statical_chart

Use this data source to get the list of DSC catalog statical chart within HuaweiCloud.

## Example Usage

```hcl
variable "label_id" {}

data "huaweicloud_dsc_catalog_statical_chart" "test" {
  label_id = var.label_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the group label ID used to filter the statical chart
  of a specific group. `label_id` and `type_id` must be specified at least one.

* `type_id` - (Optional, String) Specifies the type ID used to filter the statical chart
  of a specific type. `label_id` and `type_id` must be specified at least one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `detection_rules` - The detection rule list.

  The [detection_rules](#detection_rules_struct) structure is documented below.

* `sensitive_col_infos` - The sensitive column information list.

  The [sensitive_col_infos](#sensitive_col_infos_struct) structure is documented below.

* `total_column_number` - The total column number.

<a name="detection_rules_struct"></a>
The `detection_rules` block supports:

* `hit_number` - The hit number of the rule.

* `rule_name` - The rule name.

<a name="sensitive_col_infos_struct"></a>
The `sensitive_col_infos` block supports:

* `color_number` - The color number.

* `level_name` - The level name.

* `sensitive_number` - The sensitive number.
