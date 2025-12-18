---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_risk_config_check_rules"
description: |-
  Use this data source to get the checklist of a specified security configuration item.
---

# huaweicloud_hss_baseline_risk_config_check_rules

Use this data source to get the checklist of a specified security configuration item.

## Example Usage

```hcl
variable "check_name" {}
variable "standard" {}

data "huaweicloud_hss_baseline_risk_config_check_rules" "test" {
  check_name = var.check_name
  standard   = var.standard
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS hosts.
  If omitted, the provider-level region will be used.

* `check_name` - (Required, String) Specifies the configuration check baseline name.
  For example, **SSH**, **CentOS 7**, **Windows**.

* `standard` - (Required, String) Specifies the standard type.
  The valid values are as follows:
  + **cn_standard**: DJCP MLPS compliance standard.
  + **hw_standard**: Cloud security practice standard.

* `result_type` - (Optional, String) Specifies the result type.
  The valid values are as follows:
  + **safe**: The item passed the check.
  + **unhandled**: The item failed the check and is not ignored.
  + **ignored**: The item failed the check but is ignored.

* `check_rule_name` - (Optional, String) Specifies the check item name.
  Fuzzy match is supported.

* `severity` - (Optional, String) Specifies the risk level.
  The valid values are as follows:
  + **Security**
  + **Low**
  + **Medium**
  + **High**
  + **Critical**

* `host_id` - (Optional, String) Specifies the host ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of the check items.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `severity` - The risk level.

* `check_name` - The configuration check baseline name.

* `check_type` - The baseline type.
  The values for `check_type` and `check_name` are the same for Linux servers. For example,
  they can both be set to **SSH** or **CentOS 7**.
  For Windows servers, the values for `check_type` and `check_name` are different. For example,
  `check_type` can be set to **Windows Server 2019 R2** or **Windows Server 2016 R2**.

* `standard` - The standard type.

* `check_rule_name` - The check item (rule) name.

* `check_rule_id` - The check item ID.

* `host_num` - The number of affected servers, that is, the number of servers where
  the current baseline check is performed.

* `scan_result` - The check result.
  The valid values are as follows:
  + **pass**: The check is passed.
  + **failed**: The check is not passed.

* `status` - The check item status.
  The valid values are as follows:
  + **safe**
  + **ignored**
  + **unhandled**
  + **fixing**
  + **fix-failed**
  + **verifying**

* `enable_fix` - Whether one-click fix is supported.
  The valid values are as follows:
  + **1**: One-click fix is supported.
  + **0**: not supported

* `enable_click` - Whether the fix, ignore, and verify buttons of the check item are enabled.
  The valid values are as follows:
  + **true**: The buttons are enabled.
  + **false**: The buttons are disabled.

* `not_enable_click_description` - The reason why it cannot be clicked.

* `rule_params` - The value range of a parameter that can be configured to fix a check item.
  This information is only returned for the parameters that can be configured to fix check items.
  The [rule_params](#rule_params_struct) structure is documented below.

<a name="rule_params_struct"></a>
The `rule_params` block supports:

* `rule_param_id` - The check item parameter ID.

* `rule_desc` - The check item parameter description.

* `default_value` - The default values of check item parameters.

* `range_min` - The minimum value of check item parameters.

* `range_max` - The maximum value of check item parameters.
