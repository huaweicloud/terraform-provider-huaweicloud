---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_security_checks_details"
description: |-
  Use this data source to get the list of HSS detailed checklist of baseline inspection items within HuaweiCloud.
---

# huaweicloud_hss_baseline_security_checks_details

Use this data source to get the list of HSS detailed checklist of baseline inspection items within HuaweiCloud.

## Example Usage

```hcl
variable "support_os" {}
variable "standard" {}

data "huaweicloud_hss_baseline_security_checks_details" "test" {
  support_os = var.support_os
  standard   = var.standard
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `support_os` - (Required, String) Specifies the supported operating system of the policy.  
  The valid values are as follows:
  + **Linux**
  + **Windows**
  + **Other**

* `standard` - (Required, String) Specifies the standard type.  
  The valid values are as follows:
  + **cn_standard**: Waiting for compliance standards.
  + **hw_standard**: Cloud Security Practice Standards.
  + **cis_standard**: General safety standards.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `check_type` - (Optional, String) Specifies the configuration check (baseline) type, such as **SSH**, **CentOS 7**,
  **Windows Server 2019 R2**, **Windows Server 2016 R2**, **MySQL5-Windows**, **weakpwd**, **pwdcomplexity**.

* `tag` - (Optional, String) Specifies the tag of baseline inspection items.

* `check_rule_name` - (Optional, String) Specifies the name of configuration check (baseline) check item.

* `severity` - (Optional, String) Specifies the risk level of configuration check (baseline) check items.

* `level` - (Optional, String) Specifies the version information of configuration check (baseline) check items.

* `group_id` - (Optional, String) Specifies the policy group ID.

* `checked` - (Optional, String) Specifies whether the default is selected or not.  
  The valid values are as follows:
  + **true**: Selected.
  + **false**: Not selected.

  Defaults to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of security check policies.

* `check_details` - The checklist.

  The [check_details](#check_details_struct) structure is documented below.

<a name="check_details_struct"></a>
The `check_details` block supports:

* `key` - The unique value of inspection item.

* `check_rule_id` - The check item ID.

* `check_rule_name` - The name of inspection item (inspection rule).

* `check_rule_type` - The check if the item type is a numerical type.  
  The valid values are as follows:
  + **1**: Yes.
  + **0**: Not.

* `check_type` - The configuration check (baseline) type.

* `severity` - The risk level of configuration check (baseline) check items.

* `level` - The version information of configuration check (baseline) check items.

* `checked` - The check if the item is selected.

* `rule_params` - The customizable parameters.

  The [rule_params](#rule_params_struct) structure is documented below.

<a name="rule_params_struct"></a>
The `rule_params` block supports:

* `rule_param_id` - The check item parameter ID.

* `rule_desc` - The description of inspection item parameters.

* `default_value` - The default values for inspection item parameters.

* `range_min` - The parameter for the inspection item can take the minimum value.

* `range_max` - The maximum value that the inspection item parameter can take.
