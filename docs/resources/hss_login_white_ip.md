---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_login_white_ip"
description: |-
  Manages an HSS login white ip operation resource within HuaweiCloud.
---

# huaweicloud_hss_login_white_ip

Manages an HSS login white ip operation resource within HuaweiCloud.

-> This resource is a one-time action resource using to operation HSS login white ip. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "white_ip" {}
variable "host_id_list" {
  type = list(string)
}

resource "huaweicloud_hss_login_white_ip" "test" {
  white_ip     = var.white_ip
  host_id_list = var.host_id_list
  enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `white_ip` - (Required, String, NonUpdatable) Specifies whitelist IP or IP network segment,
  a single account can add up to `10` SSH login IP whitelist.

* `host_id_list` - (Required, List, NonUpdatable) Specifies the list of host IDs. When deleting whitelist IPs or
  IP segments, the server ID list needs to be set to an empty list.

* `enabled` - (Optional, Bool, NonUpdatable) Specifies whitelist activation status.  
  The valid values are as follows:
  + **true**: Enabled.
  + **false**: Disabled.

  Defaults to **false**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
