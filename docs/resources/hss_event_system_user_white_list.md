---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_system_user_white_list"
description: |-
  Manages an HSS event system user white list resource within HuaweiCloud.
---

# huaweicloud_hss_event_system_user_white_list

Manages an HSS event system user white list resource within HuaweiCloud.

## Example Usage

```hcl
variable "host_id" {}

resource "huaweicloud_hss_event_system_user_white_list" "test" {
  host_id               = var.host_id
  system_user_name_list = ["test_user1", "test_user2"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `host_id` - (Required, String, NonUpdatable) Specifies the host ID.

* `system_user_name_list` - (Required, List) Specifies the list of system user-names to be added to the white list.

* `remarks` - (Optional, String) Specifies the remarks of the white list.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `delete_all` - (Optional, Bool) Specifies whether to delete all system user white lists. When set to `true`, all
  system user white lists under HSS will be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID same as `host_id`.

* `enterprise_project_name` - The enterprise project name.

* `host_name` - The host name.

* `private_ip` - The private IP address of the host.

* `public_ip` - The public IP address of the host.

* `update_time` - The update time in milliseconds.

## Import

The HSS event system user white list can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_hss_event_system_user_white_list.test <host_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `enterprise_project_id`, `delete_all`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_hss_event_system_user_white_list" "test" {
  ...

  lifecycle {
    ignore_changes = [
      enterprise_project_id, delete_all,
    ]
  }
}
```
