---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_overviews"
description: |-
  Use this data source to get the overviews of baseline check.
---

# huaweicloud_hss_baseline_overviews

Use this data source to get the overviews of baseline check.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_overviews" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the policy group ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scan_time` - The last detection time, in milliseconds.

* `host_num` - The number of checked servers.

* `failed_host_num` - The number of servers that failed the check.

* `check_type_num` - The number of detection baselines.

* `check_rule_num` - The number of detection rules.

* `check_rule_pass_rate` - The pass rate of baseline inspection items.

* `cn_standard_check_rule_pass_rate` - The pass rate of the baseline inspection items for cloud security practices.

* `hw_standard_check_rule_pass_rate` - The pass rate of the baseline inspection items for classified protection
  compliance.

* `check_rule_failed_num` - The number of failed inspection items.

* `check_rule_high_risk` - The number of high-risk examination items.

* `check_rule_medium_risk` - The number of medium-risk examination items.

* `check_rule_low_risk` - The number of low-risk examination items.

* `weak_pwd_total_host` - The total number of hosts for weak password detection.

* `weak_pwd_risk` - The number of hosts with weak passwords

* `weak_pwd_normal` - The number of hosts without weak passwords.

* `weak_pwd_not_protected` - The number of hosts without protection enabled.

* `host_risks` - The TOP5 list of server risks.
  The [host_risks](#baseline_host_risks) structure is documented below.

* `weak_pwd_risk_hosts` - The TOP5 list of weak password risks for hosts
  The [weak_pwd_risk_hosts](#baseline_weak_pwd_risk_hosts) structure is documented below.

<a name="baseline_host_risks"></a>
The `host_risks` block supports:

* `host_id` - The host ID.

* `host_name` - The server name.

* `host_ip` - The server IP address.

* `scan_time` - The scan time, in milliseconds.

* `high_risk_num` - The number of high-risk risks.

* `medium_risk_num` - The number of medium-risk risks.

* `low_risk_num` - The number of low-risk risks.

<a name="baseline_weak_pwd_risk_hosts"></a>
The `weak_pwd_risk_hosts` block supports:

* `host_id` - The host ID.

* `host_name` - The server name.

* `host_ip` - The server IP address.

* `weak_pwd_num` - The number of weak passwords.
