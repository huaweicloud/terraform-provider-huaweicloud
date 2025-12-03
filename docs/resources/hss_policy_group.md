---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_policy_group"
description: |-
  Manages a policy group resource within HuaweiCloud HSS.
---

# huaweicloud_hss_policy_group

Manages a policy group resource within HuaweiCloud HSS.

## Example Usage

```hcl
variable "group_id" {}

resource "huaweicloud_hss_policy_group" "test" {
  group_id              = var.group_id
  name                  = "test-policy-group"
  description           = "Test policy group"
  protect_mode          = "equalization"
  enterprise_project_id = "all_granted_eps"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the policy group to be copied. Only premium and
  container edition policy groups can be copied.
  The value of this field can be obtained through the datasource `huaweicloud_hss_policy_groups`.
  The `group_id` whose `support_version` is **hss.version.container.enterprise** or **hss.version.premium** is the ID
  of the policy group that can be copied. This field could not be retrieved and populated during the query.

* `name` - (Required, String, NonUpdatable) Specifies the policy group name.

* `description` - (Optional, String, NonUpdatable) Specifies the description of a policy group.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the policy
  group belongs. This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

* `protect_mode` - (Optional, String) Specifies the protection mode.
  The valid values are:
  + **high_detection**: sensitive mode
  + **equalization**: balanced mode

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as the policy group ID).

* `host_num` - The number of associated servers.

* `default_group` - Whether it is the default policy group.

* `deletable` - Whether deletion is allowed.
  Deletion is allowed only when `default_group` is **false** and `host_num` is `0`.

* `support_os` - The operating systems supported by the policy group. Valid values are **Linux** and **Windows**.

* `support_version` - The versions supported by the policy group. Valid values are:
  + **hss.version.advanced**
  + **hss.version.enterprise**
  + **hss.version.premium**
  + **hss.version.wtp**
  + **hss.version.container.enterprise**

## Import

Policy group can be imported using the `enterprise_project_id` and `id` separated by a slash, e.g.

### Import resource under the default enterprise project

```bash
$ terraform import huaweicloud_hss_policy_group.test 0/<id>
```

### Import resource from non-default enterprise project

```bash
$ terraform import huaweicloud_hss_policy_group.test <enterprise_project_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `group_id`, and `enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_policy_group" "test" { 
  # ...

  lifecycle {
    ignore_changes = [
      group_id,
      enterprise_project_id,
    ]
  }
}
```
