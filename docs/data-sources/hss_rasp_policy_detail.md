---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_policy_detail"
description: |-
  Use this data source to query a specifies protection policy details.
---

# huaweicloud_hss_rasp_policy_detail

Use this data source to query a specifies protection policy details.

## Example Usage

```hcl
variable "policy_id" {}

data "huaweicloud_hss_rasp_policy_detail" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the protection policy ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `policy_name` - The protection policy name.

* `os_type` - The OS type.

* `rule_list` - All ports that match the filter parameters.
  The [rule_list](#rule_list_struct) structure is documented below.

<a name="rule_list_struct"></a>
The `rule_list` block supports:

* `chk_feature_id` - The detection feature rule ID.

* `chk_feature_name` - The detection feature rule name.

* `chk_feature_desc` - The detection feature rule description.

* `feature_configure` - The detection feature rule configuration information.

* `protective_action` - The default protection action.
  The valid values are as follows:
  + `1`: Detection.
  + `2`: Detection and blocking/interception.

* `optional_protective_action` - The available protection action.
  The valid values are as follows:
  + `1`: Detection.
  + `2`: Detection and blocking/interception.
  + `3`: All.

* `enabled` - The enabling status.
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `editable` - Whether the configuration information can be edited.
  The valid values are as follows:
  + `0`: No.
  + `1`: Yes.
