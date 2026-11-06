---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_associated_resource_rules"
description: |-
  Use this data source to query the TMS associated resource rules within HuaweiCloud.
---

# huaweicloud_tms_associated_resource_rules

Use this data source to query the TMS associated resource rules within HuaweiCloud.

## Example Usage

```hcl
# Query all rules
data "huaweicloud_tms_associated_resource_rules" "all" {}

# Query rules by region
data "huaweicloud_tms_associated_resource_rules" "filtered" {
  region_id = "cn-north-4"
}
```

## Argument Reference

The following arguments are supported:

* `setting_name` - (Optional, String) Specifies the name of the rule setting used to filter rules.

* `region_id` - (Optional, String) Specifies the region ID used to filter rules.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of associated resource rules.
  The [rules](#tms_associated_resource_rules_attr) structure is documented below.

<a name="tms_associated_resource_rules_attr"></a>
The `rules` block supports:

* `setting_name` - The name of the rule setting.

* `region_id` - The region ID where the rule takes effect.

* `tag_keys` - The list of tag keys that the rule takes effect on.

* `existing_resource_status` - Whether the rule takes effect for existing resources.
  The value can be **enabled** or **disabled**.

* `auto_delete_status` - Whether to automatically delete tags when the association is removed.
  The value can be **enabled** or **disabled**.