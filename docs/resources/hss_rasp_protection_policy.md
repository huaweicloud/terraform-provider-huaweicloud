---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_rasp_protection_policy"
description: |-
  Manages a protection policy resource within HuaweiCloud.
---

# huaweicloud_hss_rasp_protection_policy

Manages a protection policy resource within HuaweiCloud.

-> Create the protection policy resource need to meet the following conditions:
  <br/>1. The application protection is available in HSS premium, WTP, and container editions.
  <br/>2. For more details, please refer to [document](https://support.huaweicloud.com/intl/en-us/usermanual-hss2.0/hss_01_0610.html).

## Example Usage

```hcl
variable "policy_name" {}
variable "os_type" {}

resource "huaweicloud_hss_rasp_protection_policy" "test" {
  policy_name = var.policy_name
  os_type     = var.os_type

  feature_list {
    chk_feature_id    = 1
    protective_action = 1
    enabled           = 0
    feature_configure = "/guiserver/rule/create"
  }

  feature_list {
    chk_feature_id    = 2
    protective_action = 1
    enabled           = 1
    feature_configure = "/guiserver/rule/update"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `policy_name` - (Required, String) Specifies the protection policy name.
  The name must be unique.

* `feature_list` - (Required, List) Specifies the detection feature rule list.
  The [feature_list](#policy_feature_list_struct) structure is documented below.

* `os_type` - (Optional, String, NonUpdatable) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="policy_feature_list_struct"></a>
The `feature_list` block supports:

* `chk_feature_id` - (Required, Int) Specifies the detection feature rule ID.

* `protective_action` - (Required, Int) Specifies the protective action.
  The valid values are as follows:
  + **1**: Indicates detect.
  + **2**: Indicates detect and block.

* `enabled` - (Required, Int) Specifies the enabled status.
  The valid values are as follows:
  + **0**: Indicates disabled.
  + **1**: Indicates enabled.

* `feature_configure` - (Required, String) Specifies the detection feature rule configuration information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `rule_list` - The port and protocol list.
  The [rule_list](#policy_rule_list_struct) structure is documented below.

<a name="policy_rule_list_struct"></a>
The `rule_list` block supports:

* `chk_feature_id` - The detection feature rule ID.

* `chk_feature_name` - The detection rule name.

* `chk_feature_desc` - The detection rule description.

* `feature_configure` - The detection feature rule configuration information.

* `protective_action` - The protective action.

* `optional_protective_action` - The optional protection action.
  The valid values are as follows:
  + **1**: Indicates detect.
  + **2**: Indicates detect and block.
  + **3**: Indicates all.

* `enabled` - The enabled status.

* `editable` - Whether the configuration information can be edited.
  The valid values are as follows:
  + **0**: Indicates no.
  + **1**: Indicates yes.

## Import

The resource can be imported using the `enterprise_project_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_hss_rasp_protection_policy.test <enterprise_project_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `feature_list`, `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_rasp_protection_policy" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      feature_list, enterprise_project_id,
    ]
  }
}
```
