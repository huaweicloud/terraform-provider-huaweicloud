---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_policies"
description: |-
  Use this data source to get the list of protection policies.
---

# huaweicloud_hss_rasp_policies

Use this data source to get the list of protection policies.

## Example Usage

```hcl
data "huaweicloud_hss_rasp_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_name` - (Optional, String) Specifies the policy name. Supports fuzzy match.

* `os_type` - (Optional, String) Specifies the operating system type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - All policies that match the filter parameters.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `policy_id` - The policy ID.

* `policy_name` - The policy name.

* `os_type` - The operating system type.

* `host_num` - The number of associated hosts.

* `rule_name` - The names of detection rules, separated by commas(,).
