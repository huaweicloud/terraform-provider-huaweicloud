---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_policy_group_deploy"
description: |-
  Manages an HSS policy group deploy resource within HuaweiCloud.
---
# huaweicloud_hss_policy_group_deploy

Manages an HSS policy group deploy resource within HuaweiCloud.

-> This resource is only a one-time action resource for HSS policy group deploy. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "target_policy_group_id" {}
variable "enterprise_project_id" {}
variable "operate_all" {}
variable "host_id_list" {
  type = list(string)
}

resource "huaweicloud_hss_policy_group_deploy" "test" {
  target_policy_group_id = var.target_policy_group_id
  enterprise_project_id  = var.enterprise_project_id
  operate_all            = var.operate_all
  host_id_list           = var.host_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the HSS policy group resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `target_policy_group_id` - (Required, String, NonUpdatable) Specifies the ID of the policy group to be deployed.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project.
  This parameter is valid only when the enterprise project function is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  If you want to deploy the policy for all hosts under all enterprise projects,
  set this parameter to ***all_granted_eps***.

* `operate_all` - (Optional, Bool, NonUpdatable) Specifies whether to deploy the policy on all hosts.
  If the value is ***true***, not need to configure `host_id_list`.
  If the value is ***false***, configure `host_id_list`.

* `host_id_list` - (Optional, List, NonUpdatable) Specifies the ID list of servers where the policy group needs to be deployed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The source ID.
