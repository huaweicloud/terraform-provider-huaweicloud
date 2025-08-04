---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_switch_honeypot_port_policy"
description: |-
  Manages a resource to switch the dynamic port honeypot policy of a server within HuaweiCloud.
---

# huaweicloud_hss_switch_honeypot_port_policy

Manages a resource to switch the dynamic port honeypot policy of a server within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "policy" {}
variable "host_id" {}

resource "huaweicloud_hss_switch_honeypot_port_policy" "test" {
  policy_id = var.policy_id
  host_id   = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `policy_id` - (Required, String, NonUpdatable) Specifies the dynamic port honeypot policy ID.

* `host_id` - (Required, String, NonUpdatable) Specifies the host ID.
  Multiple host ID is supported, separated by commas (,).

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
