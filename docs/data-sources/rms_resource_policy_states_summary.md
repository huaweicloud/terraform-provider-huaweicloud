---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_policy_states_summary"
description: |-
  Use this data source to get the summary of a resource policy states .
---

# huaweicloud_rms_resource_policy_states_summary

Use this data source to get the summary of a resource policy states .

## Example Usage

```hcl
data "huaweicloud_rms_resource_policy_states_summary" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `resource_name` - (Optional, String) Specifies the resource name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - Indicates the value of querying the compliance result.
  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `compliance_state` - Indicates the rule status.

* `results` - Indicates the results of compliance summaries.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `resource_details` - Indicates the resource compliance summary details.
  The [resource_details](#results_details_struct) structure is documented below.

* `assignment_details` - Indicates the compliance summary details.
  The [assignment_details](#results_details_struct) structure is documented below.

<a name="results_details_struct"></a>
The `resource_details` and `assignment_details` block supports:

* `compliant_count` - Indicates the number of compliant resources.

* `non_compliant_count` - Indicates the number of non-compliant resources.
