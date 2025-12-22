---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_check_rule_fail_detail"
description: |-
  Use this data source to get the list of HSS baseline check rule fail detail within HuaweiCloud.
---

# huaweicloud_hss_baseline_check_rule_fail_detail

Use this data source to get the list of HSS baseline check rule fail detail within HuaweiCloud.

## Example Usage

```hcl
variable "check_rule_id" {}

data "huaweicloud_hss_baseline_check_rule_fail_detail" "test" {
  check_rule_id = var.check_rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `check_rule_id` - (Required, String) Specifies the check item ID.
  The value of this parameter can query from dataSource `huaweicloud_hss_baseline_risk_config_check_rules`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_id` - (Optional, String) Specifies the server ID.

* `check_name` - (Optional, String) Specifies the configuration check (baseline) name.
  For example: **SSH**, **CentOS 7**, **Windows**.

* `standard` - (Optional, String) Specifies the standard type. Valid values are:
  + **cn_standard**: DJCP MLPS compliance standard.
  + **hw_standard**: Cloud security practice standard.
  + **cis_standard**: Common security standard.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `fail_detail_list` - The list of check rule fix fail detail.

  The [fail_detail_list](#fail_detail_list_struct) structure is documented below.

<a name="fail_detail_list_struct"></a>
The `fail_detail_list` block supports:

* `fix_fail_reason` - The reason for fix failure.

* `host_name` - The server name.
