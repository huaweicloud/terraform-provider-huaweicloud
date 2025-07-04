---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_vulnerabilities"
description: |-
  Use this data source to get the list of vulnerabilities.
---

# huaweicloud_hss_vulnerabilities

Use this data source to get the list of vulnerabilities.

## Example Usage

```hcl
data "huaweicloud_hss_vulnerabilities" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vul_id` - (Optional, String) Specifies the vulnerability ID.

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

* `cve_id` - (Optional, String) Specifies the CVE ID.

* `label_list` - (Optional, String) Specifies the vulnerability tags.
  Multiple labels can be transferred for filtering, separated by commas (,). For example, **test1,test2**.

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**
  + **common**
  + **test**

* `group_name` - (Optional, String) Specifies the server group name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - Software vulnerability list

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `vul_id` - The vulnerability ID.

* `vul_name` - The vulnerability name.

* `type` - The vulnerability type.

* `repair_necessity` - The necessity of fixing a vulnerability.
  The valid values are as follows:
  + **Critical**: The CVSS score of the vulnerability is greater than or equal to `9`, corresponding to the high risk
  level on the console.
  + **High**: The CVSS score of the vulnerability is greater than or equal to `7` and less than `9`, corresponding to
  the medium risk level on the console.
  + **Medium**: The CVSS score of the vulnerability is greater than or equal to `4` and less than `7`, corresponding to
  the medium risk level on the console.
  + **Low**: The CVSS score of the vulnerability is less than `4`, corresponding to the low risk level on the console.

* `severity_level` - The vulnerability severity.
  The valid values are as follows:
  + **Critical**: The CVSS score of the vulnerability is greater than or equal to `9`, corresponding to the high risk
  level on the console.
  + **High**: The CVSS score of the vulnerability is greater than or equal to `7` and less than `9`, corresponding to
  the medium risk level on the console.
  + **Medium**: The CVSS score of the vulnerability is greater than or equal to `4` and less than `7`, corresponding to
  the medium risk level on the console.
  + **Low**: The CVSS score of the vulnerability is less than 4, corresponding to the low risk level on the console.

* `description` - The vulnerability description

* `label_list` - The vulnerability tags list.

* `host_num` - The number of affected servers.

* `unhandle_host_num` - The number of unhandled servers, excluding ignored and fixed servers.

* `solution_detail` - The vulnerability fixing guide.

* `url` - The vulnerability URL.

* `host_id_list` - The list of servers where the vulnerability can be handled.

* `cve_list` - The CVE list

  The [cve_list](#data_list_cve_list_struct) structure is documented below.

* `patch_url` - The patch address.

* `repair_priority` - The fix priority.

* `hosts_num` - The number of affected servers.

  The [hosts_num](#data_list_hosts_num_struct) structure is documented below.

* `repair_success_num` - The number of successful repairs.

* `fixed_num` - The number of repairs.

* `ignored_num` - The number of ignored.

* `verify_num` - The number of verifications.

* `repair_priority_list` - The number of servers corresponding to each fixing priority.

  The [repair_priority_list](#data_list_repair_priority_list_struct) structure is documented below.

* `scan_time` - The latest scan time.

<a name="data_list_cve_list_struct"></a>
The `cve_list` block supports:

* `cve_id` - The CVE ID.

* `cvss` - The CVSS score.

<a name="data_list_hosts_num_struct"></a>
The `hosts_num` block supports:

* `common` - The number of common servers.

* `test` - The number of test servers.

* `important` - The number of important servers.

<a name="data_list_repair_priority_list_struct"></a>
The `repair_priority_list` block supports:

* `repair_priority` - The fixing priority.

* `host_num` - The number of servers corresponding to the fixing priority.
