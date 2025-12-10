---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_overview_security_risks"
description: |-
  Use this data source to get the list of security risks.
---

# huaweicloud_hss_overview_security_risks

Use this data source to get the list of security risks.

## Example Usage

```hcl
data "huaweicloud_hss_overview_security_risks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_risk` - The intrusion risks.
  The [alarm_risk](#alarm_risk_struct) structure is documented below.

* `baseline_risk` - The baseline risks.
  The [baseline_risk](#baseline_risk_struct) structure is documented below.

* `asset_risk` - The asset risks.
  The [asset_risk](#asset_risk_struct) structure is documented below.

* `security_protect_risk` - The security protection risks.
  The [security_protect_risk](#security_protect_risk_struct) structure is documented below.

* `vul_risk` - The vulnerability risks.
  The [vul_risk](#vul_risk_struct) structure is documented below.

* `image_risk` - The list of hosts affected by the vulnerability.
  The [image_risk](#image_risk_struct) structure is documented below.

* `total_risk_num` - The total number of risks.

<a name="alarm_risk_struct"></a>
The `alarm_risk` block supports:

* `risk_list` - The list of risks.
  The [risk_list](#risk_list_info_struct) structure is documented below.

* `deduct_score` - The score deduction.

* `policy_list` - The policy information.
  The [policy_list](#policy_list_info_struct) structure is documented below.

* `total_risk_num` - The total number of risks.

<a name="baseline_risk_struct"></a>
The `baseline_risk` block supports:

* `risk_list` - The list of baseline risks.
  The [risk_list](#risk_list_info_struct) structure is documented below.

* `deduct_score` - The score deduction.

* `policy_list` - The list of policies that are not enabled.
  The [policy_list](#policy_list_info_struct) structure is documented below.

* `existed_pwd_host_num` - The number of servers with weak passwords.

* `un_scanned_baseline_host_num` - The number of servers where baseline check is not performed.

* `total_risk_num` - The total number of risks.

<a name="asset_risk_struct"></a>
The `asset_risk` block supports:

* `existed_danger_port_host_num` - The number of servers with dangerous ports.

* `policy_list` - The list of policies.
  The [policy_list](#policy_list_info_struct) structure is documented below.

* `deduct_score` - The score deduction.

* `total_risk_num` - The total number of risks.

<a name="security_protect_risk_struct"></a>
The `security_protect_risk` block supports:

* `un_open_protection_host_num` - The number of unprotected servers.

* `deduct_score` - The score deduction.

* `total_risk_num` - The total number of risks.

<a name="vul_risk_struct"></a>
The `vul_risk` block supports:

* `risk_list` - The list of vulnerability risks.
  The [risk_list](#risk_list_info_struct) structure is documented below.

* `deduct_score` - The score deduction.

* `un_scanned_host_num` - The number of servers where vulnerability scan is not performed (in the past month).

* `total_risk_num` - The total number of risks.

<a name="image_risk_struct"></a>
The `image_risk` block supports:

* `deduct_score` - The score deduction.

* `un_scanned_image_num` - The number of unscanned images.

* `risk_list` - The list of images risks.
  The [risk_list](#risk_list_image_info_struct) structure is documented below.

* `total_risk_num` - The total number of risks.

<a name="risk_list_info_struct"></a>
The `risk_list` block supports:

* `severity` - The risk level.

* `risk_num` - The number of risks.

* `effected_host_num` - The number of affected assets.

<a name="policy_list_info_struct"></a>
The `policy_list` block supports:

* `policy_id` - The policy ID.

* `policy_name` - The policy name.

* `os_type` - The OS type.

* `host_num` - The associated servers.

* `rule_name` - The detection feature rule name.

<a name="risk_list_image_info_struct"></a>
The `risk_list` block supports:

* `severity` - The risk level.

* `image_num` - The number of images.
