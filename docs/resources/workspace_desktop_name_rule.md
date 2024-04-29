---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_name_rule"
description: ""
---

# huaweicloud_workspace_desktop_name_rule

Using this resource to manage a name rule for the desktops creation within HuaweiCloud.

## Example Usage

```hcl
variable "policy_name" {}

resource "huaweicloud_workspace_desktop_name_rule" "test" {
  name                         = var.policy_name
  name_prefix                  = "test$DomainUser$end"
  digit_number                 = 3
  start_number                 = 2
  single_domain_user_increment = 1
  is_default_policy            = true 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the rule.
  The name can contain `1` to `30` characters, only digits, letters and underscores (_) are allowed.
  The name must start with a letter or an underscore.
  
* `name_prefix` - (Required, String) Specifies the prefix of desktop name.
  The format are as follows: `$DomainUser$`, `xx$DomainUser$`, `$DomainUser$xx`, `xx$DomainUser$xx`, `xx`.
  The value only letters, digits, hyphens (-) and `$DomainUser$` are allowed, it must start with a letter, digit or `$DomainUser$`.
  If the value contains `$DomainUser$` string, it means that the rule contains the user name. Otherwise, it means that
  the rule doesn't contain user name.
  The former format can't be used for desktop pool, and the latter can be used in all scenarios.

  -> The desktop name format consists of `name_prefix` and `digit_number`, and the total length cannot exceed 15 characters.
  The desktop name format are as follows:
  Include user name: `A+username+B+1`.
  Without user name: `A+1`.

* `digit_number` - (Required, Int) Specifies the number of valid digits in the desktop name suffix.
  The valid value is range from `1` to `5`.

* `start_number` - (Required, Int) Specifies the start number of the desktop name suffix.
  The value is related to the `digit_number` parameter.
  If the `digit_number` parameter is set to `1`, the valid value is range from `1` to `9`.
  If the `digit_number` parameter is set to `2`, the valid value is range from `1` to `99`, and so on.

* `single_domain_user_increment` - (Required, Int) Specifies whether to increment by single user name.
  The valid values are as follows:
  + **1**: Increment by single user name.
  + **0**: Increment by tenant.

  e.g. Assume that there are three user A, B, C. Allocate two desktops to user A, one to user B, and one to user C.
  The `digit_number` parameter is set to `2`, the `start_number` parameter is set to `1`.
  If the `single_domain_user_increment` is set to `1`, the desktop names are as follows: `A01`, `A02`, `B01`, `C01`.
  If the `single_domain_user_increment` is set to `0`, the desktop names are as follows: `A01`, `A02`, `B03`, `C04`.

  -> If the `name_prefix` parameter not contain `$DomainUser$`, the `single_domain_user_increment` value must set to `0`.

* `is_default_policy` - (Optional, Bool) Specifies whether to set as default rule. The default value is **false**.
  
  ~> Only one default policy is allowed. If the current resource is set as the default policy, the value of
     `is_default_policy` of the original default rule will be changed to `false`. If that rule is also managed by the provider,
     please modify the corresponding script simultaneously, otherwise it will cause changes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_contain_user` - Whether the desktop name contains the user name.

## Import

The desktop name rule can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_desktop_name_rule.test <id>
```
