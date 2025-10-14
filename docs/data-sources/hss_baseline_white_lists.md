---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_white_lists"
description: |-
  Use this data source to get the list of baseline white lists within HuaweiCloud.
---

# huaweicloud_hss_baseline_white_lists

Use this data source to get the list of baseline white lists within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_white_lists" "test" {
  enterprise_project_id = "0"
  os_type               = "Linux"
  rule_type             = "all_host"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `check_rule_name` - (Optional, String) Specifies the baseline check rule name to filter the white lists.

* `os_type` - (Optional, String) Specifies the operating system type to filter the white lists.
  Valid values are **Linux** and **Windows**.

* `rule_type` - (Optional, String) Specifies the rule type of the white lists.
  Valid values are **specific_host** and **all_host**.

* `tag` - (Optional, String) Specifies the tag to filter the white lists.

* `description` - (Optional, String) Specifies the description to filter the white lists.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of baseline white lists.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The ID of the baseline white list.

* `rule_type` - The rule type of the baseline white list.

* `os_type` - The operating system type of the baseline white list.

* `index_version` - The index version of the baseline check rule.

* `check_type` - The check type of the baseline white list.

* `standard` - The standard type of the baseline white list.
  Valid values are **cn_standard**, **hw_standard**, and **cis_standard**.

* `tag` - The tag of the baseline white list.

* `check_rule_name` - The name of the baseline check rule.

* `description` - The description of the baseline white list.
