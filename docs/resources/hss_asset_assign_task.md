---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_assign_task"
description: |-
  Manages a resource to create a global asset scan task within HuaweiCloud.
---

# huaweicloud_hss_asset_assign_task

Manages a resource to create a global asset scan task within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "category" {}

resource "huaweicloud_hss_asset_assign_task" "test" {
  category  = var.category
  all_hosts = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `category` - (Required, String, NonUpdatable) Specifies the asset category.
  The valid values are as follows:
  + **host**
  + **container**

* `all_hosts` - (Required, Bool, NonUpdatable) Specifies Whether scan all hosts.
  The valid values are as follows:
  + **true**
  + **false**

* `host_ids` - (Optional, List, NonUpdatable) Specifies the host ID list.

  -> This parameter is valid and required when `all_hosts` is set to **false**.

* `server_group` - (Optional, List, NonUpdatable) Specifies the host group ID list.

  -> This parameter is valid when `all_hosts` is set to **false**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `error` - The reasons for failure.
