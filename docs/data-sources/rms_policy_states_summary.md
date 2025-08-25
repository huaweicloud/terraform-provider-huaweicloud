---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_states_summary"
description: |-
  Use this data source to get compliance summary by user.
---

# huaweicloud_rms_policy_states_summary

Use this data source to get compliance summary by user.

## Example Usage

```hcl
data "huaweicloud_rms_policy_states_summary" "test" {}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Optional, List) Specifies the tag list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - Indicates the results of compliance summaries.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `resource_details` - Indicates the resource compliance summary details.
  The [resource_details](#details_struct) structure is documented below.

* `assignment_details` - Indicates the compliance summary details.
  The [assignment_details](#details_struct) structure is documented below.

* `group_name` - Indicates the group name

* `group_account_name` - Indicates the acccount name

<a name="details_struct"></a>
The `resource_details` and `assignment_details` block supports:

* `compliant_count` - Indicates the number of compliant resources.

* `non_compliant_count` - Indicates the number of non-compliant resources.
