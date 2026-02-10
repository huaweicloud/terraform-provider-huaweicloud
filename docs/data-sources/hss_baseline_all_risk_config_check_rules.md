---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_all_risk_config_check_rules"
description: |-
  Use this data source to get the list of HSS baseline all risk config check rules within HuaweiCloud.
---

# huaweicloud_hss_baseline_all_risk_config_check_rules

Use this data source to get the list of HSS baseline all risk config check rules within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_all_risk_config_check_rules" "test" {
  type       = "linux"
  image_type = "private_image"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `check_type` - (Optional, String) Specifies the configuration check (baseline) type, such as **SSH**, **CentOS 7**,
  **Windows**.

* `standard` - (Optional, String) Specifies the baseline standard type.  
  The valid values are as follows:
  + **cn_standard**: Waiting for compliance standards.
  + **hw_standard**: Cloud Security Practice Standards.
  + **cis_standard**: General safety standards.

* `statistics_scan_result` - (Optional, String) Specifies the type of statistical result.  
  The valid values are as follows:
  + **pass**: Passed, indicating that all inspection items for the host have been checked and passed.
  + **failed**: Failed, indicating that all check items on the host have not passed and have not been processed.
  + **processed**: Processed, indicating that there are check items on the host that have not passed or have all been
    processed (ignored, highlighted).

* `check_rule_name` - (Optional, String) Specifies the name of the inspection item (inspection rule) and support fuzzy
  matching.

* `severity` - (Optional, String) Specifies the risk level.  
  The valid values are as follows:
  + **Security**: Security.
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.
  + **Critical**: Critical.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `tag` - (Optional, String) Specifies the type of baseline inspection item.

* `policy_group_id` - (Optional, String) Specifies the policy group ID. When not set, check all hosts of the tenant.
  If `host_id` exists, this value is invalid.

* `statistics_flag` - (Optional, Bool) Specifies whether to display data from statistical dimensions.  
  The valid values are **true** and **false**, defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of check rules.

* `pass_num` - The number of checked items passed.

* `failed_num` - The number of failed inspection items.

* `processed_num` - The number of processed inspection items.

* `data_list` - The data list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `tag` - The type of baseline inspection item.

* `check_rule_name` - The name of inspection item (inspection rule).

* `check_rule_id` - The check item ID.

* `severity` - The risk level.

* `check_type` - The type of configuration check (baseline).

* `check_type_name` - The type name of configuration check (baseline).

* `standard` - The baseline standard type.

* `host_num` - The number of affected servers, the number of servers undergoing current baseline detection.

* `failed_num` - This test item failed, with the number of hosts that were neither ignored nor whitelisted.

* `scan_time` - The latest detection time (ms).  
  The valid values are as follows:
  + **1**: One-click fix is supported.
  + **0**: not supported

* `statistics_scan_result` - The type of statistical result.

* `enable_fix` - Whether one-click fix is supported.  
  The valid values are as follows:
  + **1**: One-click fix is supported.
  + **0**: not supported

* `enable_click` - Whether the fix, ignore, and verify buttons of the check item are enabled.  
  The valid values are as follows:
  + **true**: The buttons are enabled.
  + **false**: The buttons are disabled.

* `cancel_ignore_enable_click` - The neglected whether the check item is clickable.  
  **true** means the button is clickable, **false** means the button is not clickable.

* `rule_params` - The value range of a parameter that can be configured to fix a check item. This information is only
  returned for the parameters that can be configured to fix check items.

  The [rule_params](#rule_params_struct) structure is documented below.

<a name="rule_params_struct"></a>
The `rule_params` block supports:

* `rule_param_id` - The check item parameter ID.

* `rule_desc` - The check item parameter description.

* `default_value` - The default values of check item parameters.

* `range_min` - The minimum value of check item parameters.

* `range_max` - The maximum value of check item parameters.
