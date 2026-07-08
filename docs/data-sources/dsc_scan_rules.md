---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_rules"
description: |-
  Use this data source to get the list of scan rules within HuaweiCloud.
---

# huaweicloud_dsc_scan_rules

Use this data source to get the list of scan rules within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_scan_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scan_rules_list` - The scan rule list.

  The [scan_rules_list](#scan_rules_list_struct) structure is documented below.

<a name="scan_rules_list_struct"></a>
The `scan_rules_list` block supports:

* `category` - The rule category.

* `project_id` - The project ID.

* `rule_desc` - The rule description.

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `rule_type` - The rule type.
