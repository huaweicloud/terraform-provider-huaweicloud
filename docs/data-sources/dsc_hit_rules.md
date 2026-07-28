---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_hit_rules"
description: |-
  Use this data source to get the list of DSC hit rules within HuaweiCloud.
---

# huaweicloud_dsc_hit_rules

Use this data source to get the list of DSC hit rules within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_dsc_hit_rules" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the scan job.

* `keyword` - (Optional, String) Specifies the sensitive object name keyword for fuzzy query.

* `asset_type` - (Optional, String) Specifies the asset type for filtering.

* `asset_id` - (Optional, String) Specifies the asset ID for filtering.

* `security_level_ids` - (Optional, List) Specifies the list of security level IDs for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `hit_rules` - The hit rule information list.

  The [hit_rules](#hit_rules_struct) structure is documented below.

<a name="hit_rules_struct"></a>
The `hit_rules` block supports:

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `top_objects` - The list of hit objects.
