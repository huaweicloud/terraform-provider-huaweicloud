---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_vulnerabilities"
description: |-
  Use this data source to get the list of vulnerabilities for a specified HSS host.
---

# huaweicloud_hss_host_vulnerabilities

Use this data source to get the list of vulnerabilities for a specified HSS host.

## Example Usage

```hcl
variable host_id {}

data "huaweicloud_hss_host_vulnerabilities" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `vul_name` - (Optional, String) Specifies the vulnerability name.

* `type` - (Optional, String) Specifies the vulnerability type.
  The valid values are as follows:
  + **linux_vul**: Linux vulnerability.
  + **windows_vul**: Windows vulnerability.
  + **web_cms**: Web-CMS vulnerability.
  + **app_vul**: Application vulnerability.
  + **urgent_vul**: Emergency vulnerability.

* `status` - (Optional, String) Specifies the vulnerability status.
  The valid values are as follows:
  + **vul_status_unfix**: Indicates not fixed.
  + **vul_status_ignored**: Indicates ignored.
  + **vul_status_verified**: Indicates verification in progress.
  + **vul_status_fixing**: Indicates fixing is in progress.
  + **vul_status_fixed**: Indicates fix succeeded.
  + **vul_status_reboot**: Indicates the issue is fixed and waiting for restart.
  + **vul_status_failed**: Indicates the issue failed to be fixed.
  + **vul_status_fix_after_reboot**: Indicates restart the host and try again.

* `handle_status` - (Optional, String) Specifies the handling status.
  The valid values are as follows:
  + **unhandled**
  + **handled**

* `repair_priority` - (Optional, String) Specifies the fixing priority.
  The valid values are as follows:
  + **Critical**
  + **High**
  + **Medium**
  + **Low**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of host vulnerabilities.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `vul_id` - The vulnerability ID.

* `vul_name` - The vulnerability name.

* `description` - The vulnerability description.

* `type` - The vulnerability type.

* `status` - The vulnerability status.

* `repair_necessity` - The repair urgency.

* `severity_level` - The risk level.
  The valid values are as follows:
  + **Critical**: The CVSS score of the vulnerability is greater than or equal to `9`, corresponding to the high risk
  level on the console.
  + **High**: The CVSS score of the vulnerability is greater than or equal to `7` and less than `9`, corresponding to
  the medium risk level on the console.
  + **Medium**: The CVSS score of the vulnerability is greater than or equal to `4` and less than `7`, corresponding to
  the medium risk level on the console.
  + **Low**: The CVSS score of the vulnerability is less than `4`, corresponding to the low risk level on the console.

* `repair_priority` - The fixing priority.

* `repair_cmd` - The repair command.

* `repair_success_num` - The total times that the vulnerability is fixed by HSS on the entire network.

* `version` - The host quota.

* `support_restore` - Whether data can be rolled back to the backup created when the vulnerability was fixed.

* `app_list` - The list of softwares affected by the vulnerability on the host.

  The [app_list](#data_list_app_list_struct) structure is documented below.

* `solution_detail` - The solution of fixed vulnerability.

* `url` - The URL.

* `app_name` - The software name.

* `app_path` - The software path.

* `app_version` - The software version.

* `is_affect_business` - Whether services are affected.

* `label_list` - The vulnerability tags list.

* `cve_list` - The CVE list.

  The [cve_list](#data_list_cve_list_struct) structure is documented below.

* `disabled_operate_types` - The list of operation types of vulnerabilities that cannot be performed.

  The [disabled_operate_types](#data_list_disabled_operate_types_struct) structure is documented below.

* `first_scan_time` - The first scan time.

* `scan_time` - The latest scan time.

<a name="data_list_app_list_struct"></a>
The `app_list` block supports:

* `app_name` - The software name.

* `app_path` - The path of the application software.
  This field is available only for application vulnerabilities.

* `app_version` - The software version.

* `upgrade_version` - The version that needs to be upgraded to fix vulnerability software.

<a name="data_list_cve_list_struct"></a>
The `cve_list` block supports:

* `cve_id` - The CVE ID.

* `cvss` - The CVSS score.

<a name="data_list_disabled_operate_types_struct"></a>
The `disabled_operate_types` block supports:

* `operate_type` - The operation type.
  The valid values are as follows:
  + **ignore**
  + **not_ignore**
  + **immediate_repair**
  + **manual_repair**
  + **verify**
  + **add_to_whitelist**

* `reason` - The reason why the operation cannot be performed.
