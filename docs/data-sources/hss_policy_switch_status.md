---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_policy_switch_status"
description: |-
  Use this data source to get the policy switch status of HSS within HuaweiCloud.
---

# huaweicloud_hss_policy_switch_status

Use this data source to get the policy switch status of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_policy_switch_status" "test" {
  enterprise_project_id = "all_granted_eps"
  policy_name           = "sp_feature"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Required, String) Specifies the enterprise project ID.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

* `policy_name` - (Required, String) Specifies the policy name. Valid value is **sp_feature**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enable` - Whether policy is switch on.
