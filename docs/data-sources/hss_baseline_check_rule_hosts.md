---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_check_rule_hosts"
description: |-
  Use this data source to get the list of hosts affected by a specified baseline check rule within HuaweiCloud.
---

# huaweicloud_hss_baseline_check_rule_hosts

Use this data source to get the list of hosts affected by a specified baseline check rule within HuaweiCloud.

## Example Usage

```hcl
variable "check_rule_id" {}
variable "standard" {}

data "huaweicloud_hss_baseline_check_rule_hosts" "test" {
  check_rule_id = var.check_rule_id
  standard      = var.standard
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `check_rule_id` - (Required, String) Specifies the baseline check rule ID.
  The value can be obtained from the response of the dataSource `huaweicloud_hss_baseline_risk_config_check_rules`.

* `standard` - (Required, String) Specifies the baseline standard type.  
  The valid values are as follows:
  + **cn_standard**: Waiting for compliance standards.
  + **hw_standard**: Cloud Security Practice Standards.
  + **qt_standard**: General safety standards.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `check_name` - (Optional, String) Specifies the configuration check (baseline) name, such as **SSH**, **CentOS 7**,
  **Windows**. Compared with `check_type`, it includes more process information such as -PID. When querying by specific
  baseline dimension, use `check_name`.

* `check_type` - (Optional, String) Specifies the configuration check (baseline) type, such as **SSH**, **CentOS 7**,
  **Windows**. When querying by check item dimension, use `check_type`.

* `result_type` - (Optional, String) Specifies the detection result type.  
  The valid values are as follows:
  + **safe**: The item passed the check.
  + **unhandled**: The item failed the check and is not ignored.
  + **ignored**: The item failed the check but is ignored.
  + **fixing**: The item is being fixed.
  + **fix-failed**: The fix failed.
  + **verifying**: The item is being verified.
  + **add_to_whitelist**: The item has been added to whitelist.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `host_name` - (Optional, String) Specifies the host name or IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of hosts affected by the baseline check rule.

* `data_list` - The list of hosts affected by the baseline check rule.
  
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `check_name` - The check name.

* `baseline_name` - The baseline name.

* `host_public_ip` - The host public IP.

* `host_private_ip` - The private IP.

* `scan_time` - The scanning time (ms).

* `failed_num` - The number of risk items.

* `passed_num` - The number of items passed.

* `diff_description` - The differentiated display prompt information.

* `description` - The ignore or add white remarks.

* `host_type` - The host type, when the host is of CCE type, return CCE.

* `enable_fix` - Does it support one click repair. `1` represents support, `0` represents no support.

* `enable_verify` - Is this check item verifiable. It requires Linux and agent version>=3.2.24.
  **true** means verifiable, **false** means unverifiable.

* `enable_click` - Is the repair, ignore, and verify button for this check item clickable.
  **true** means the button is clickable, **false** means the button is not clickable.

* `cancel_ignore_enable_click` - The neglected whether the check item is clickable.
  **true** means the button is clickable, **false** means the button is not clickable.

* `result_type` - The detection result type.

* `fix_failed_reason` - The reason for repair failure.

* `cluster_id` - The cluster ID.
