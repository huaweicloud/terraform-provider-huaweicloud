---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_event_unblock_ip"
description: |-
  Manages an HSS unblock IP resource within HuaweiCloud.
---

# huaweicloud_hss_event_unblock_ip

Manages an HSS unblock IP resource within HuaweiCloud.

-> This resource is only a one-time action resource for HSS unblock IP. Deleting this resource will not clear the
  corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "host_id" {}
variable "src_ip" {}
variable "login_type" {}

resource "huaweicloud_hss_event_unblock_ip" "test" {
  data_list {
    host_id    = var.host_id
    src_ip     = var.src_ip
    login_type = var.login_type
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the resource belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `data_list` - (Required, List, NonUpdatable) Specifies the IP list that needs to be unblocked.  
  The [data_list](#data_list_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the hosts
  belong.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - (Required, String, NonUpdatable) Specifies the host ID.

* `src_ip` - (Required, String, NonUpdatable) Specifies the IP address of the attack source.

* `login_type` - (Required, String, NonUpdatable) Specifies login type.  
  The valid values are as follows:
  + **mysql**: Represents the MySQL service.
  + **rdp**: Represents the RDP service.
  + **ssh**: Represents the SSH service.
  + **vsftp**: Represents the VSFTP service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
