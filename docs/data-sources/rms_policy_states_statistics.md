---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_states_statistics"
description: |-
  Use this data source to get the current account policy statistics.
---

# huaweicloud_rms_policy_states_statistics

Use this data source to get the current account policy statistics.

## Example Usage

```hcl
data "huaweicloud_rms_policy_states_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - Indicates the policy states statistic.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `total_resource_count` - Indicates the total count of resources.

* `non_compliant_resource_count` - Indicates  the count of noncompliance resources.

* `total_policy_count` - Indicates the total count of policy assignments.

* `non_compliant_policy_count` - Indicates the count of noncompliance policy assignments

* `statistic_date` - Indicates  the statistic date.
