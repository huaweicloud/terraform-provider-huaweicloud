---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_change_host_ignore_status"
description: |-
  Manages an HSS change host ignore status operation resource within HuaweiCloud.
---

# huaweicloud_hss_change_host_ignore_status

Manages an HSS change host ignore status operation resource within HuaweiCloud.

-> This resource is a one-time action resource using to operation HSS change host ignore status. Deleting this resource
  will not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "host_id_list" {
  type = list(string)
}

resource "huaweicloud_hss_change_host_ignore_status" "test" {
  operate_type = "ignore"
  host_id_list = var.host_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `operate_type` - (Required, String, NonUpdatable) Specifies the operation type.  
  The valid values are as follows:
  + **ignore**: Ignore host.
  + **un_ignore**: Cancel ignore host.

* `host_id_list` - (Required, List, NonUpdatable) Specifies the list of host IDs.  
  -> The host in protection cannot be ignored.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
