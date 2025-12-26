---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group_scaling_policy"
description: |-
  Manages a Workspace APP server group scaling policy resource within HuaweiCloud.
---

# huaweicloud_workspace_app_server_group_scaling_policy

Manages a Workspace APP server group scaling policy resource within HuaweiCloud.

-> An server group can only have one scaling policy resource.

## Example Usage

```hcl
variable "server_group_id" {}

resource "huaweicloud_workspace_app_server_group_scaling_policy" "test" {
  server_group_id        = var.server_group_id
  max_scaling_amount     = 10
  single_expansion_count = 1

  scaling_policy_by_session {
    session_usage_threshold           = 80
    shrink_after_session_idle_minutes = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the scaling policy is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_group_id` - (Required, String, NonUpdatable) Specifies the ID of the server group to which the scaling
  policy belongs.

* `max_scaling_amount` - (Required, Int) Specifies the maximum number of instances that can be scaled out.
  The valid value is range from `1` to `100`.

* `single_expansion_count` - (Required, Int) Specifies the number of instances to scale out in a single scaling operation.
  The valid value is range from `1` to `10`.

* `scaling_policy_by_session` - (Required, List) Specifies the session-based scaling policy configuration.
  The [scaling_policy_by_session](#workspace_scaling_policy_by_session_object) structure is documented below.

<a name="workspace_scaling_policy_by_session_object"></a>
The `scaling_policy_by_session` block supports:

* `session_usage_threshold` - (Required, Int) Specifies the total session usage threshold of the server group.

* `shrink_after_session_idle_minutes` - (Required, Int) Specifies the number of minutes to wait before releasing instances
  with no session connections.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `enable` - Whether the scaling policy is enabled.

## Import

The APP server group scaling policy can be imported using the server group ID, e.g.

```bash
$ terraform import huaweicloud_workspace_app_server_group_scaling_policy.test <id>
```
