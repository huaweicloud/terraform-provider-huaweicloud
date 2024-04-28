---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktops"
description: ""
---

# huaweicloud_workspace_desktops

Use this data source to get the list of Workspace desktops within HuaweiCloud.

## Example Usage

```hcl
variable "desktop_id" {}

data "huaweicloud_workspace_desktops" "test" {
  desktop_id = var.desktop_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `desktop_id` - (Optional, String) Specifies ID of the desktop.

* `name` - (Optional, String) Specifies the name of the desktop.

* `user_name` - (Optional, String) Specifies the user name to which the desktops belongs.

* `fixed_ip` - (Optional, String) Specifies the fixed IP address of the desktop.

* `desktop_type` - (Optional, String) Specifies the type of the desktops.
  The valid values are as follows:
  + **DEDICATED**: Normal desktop.
  + **POOLED**: Desktop in the Workspace desktop pool.

* `tags` - (Optional, Map) Specifies the key/value pairs used to query the desktops.

* `user_attached` - (Optional, String) Specify whether to query desktops by assigned users.
  The value can be **true** and **false**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the desktops.

* `image_id` - (Optional, String) Specifies the image ID of the desktops.

* `in_maintenance_mode` - (Optional, String) Specify whether to query desktops by maintenance mode.
  The value can be **true** and **false**.

* `status` - (Optional, String) Specifies the status of the desktops.
  The valid values are as follows: **ACTIVE**, **SHUTOFF**, **ERROR**.

* `subnet_id` - (Optional, String) Specifies the subnet ID of desktops.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `desktops` - The list of the desktops.  
  The [desktops](#workspace_desktops) structure is documented below.

<a name="workspace_desktops"></a>
The `desktops` block supports:

* `id` - The desktop ID.

* `name` - The desktop name.

* `type` - The desktop type.

* `enterprise_project_id` - The enterprise project ID to which the desktop belongs.

* `image_id` - The image ID of the desktop.

* `status` - The status of the desktop.

* `in_maintenance_mode` - Whether the desktop is in maintenance mode.

* `subnet_id` - The subnet ID to which the desktop belongs.

* `tags` - The key/value pairs to associate with the desktop.

* `internet_mode` - The access mode of desktop.
  The valid values are as follows: **NAT**, **EIP**, **BOTH**.

* `ip_addresses` - The list of fixed IP addresses.

* `root_volume` - The root volume to which the desktop belongs.
  The [root_volume](#desktop_volume) structure is documented below.

* `data_volume` - The data volumes to which the desktop belongs.
  The [data_volume](#desktop_volume) structure is documented below.

* `availability_zone` - The availability zone to which the desktop belongs.

* `attach_user_infos` - The list of user information assigned to the desktop.
  The [attach_user_infos](#desktop_attach_user_infos) structure is documented below.

* `created_at` - The creation time of the desktop.

* `site_name` - The site name of desktop.

* `site_type` - The site type of desktop.

* `attach_state` - The assignment status of the desktop.

* `product_id` - The product ID used by the desktop.

* `flavor_id` - The flavor ID used by the desktop.

* `ou_name` - The ou name used by the desktop.

* `ou_version` - The ou version used by the desktop.

* `join_domain` - Whether Connect to AD domain when the desktop was created. The valid values are as follows:
  + **0**: Indicates AD domain enabled.
  + **1**: Indicates AD domain not enabled.

* `is_support_internet` - Whether the desktop is associated with an EIP.

<a name="desktop_attach_user_infos"></a>
The `attach_user_infos` block supports:

* `user_group` - The user group to which the desktop belongs.

* `user_name` - The user name to which the desktop belongs.

<a name="desktop_volume"></a>
The `root_volume` and `data_volume` block supports:

* `type` - The volume type.

* `size` - The volume size.

* `device` - The device of the volume.

* `id` - The unique identification ID of the disk.

* `name` - The volume name.

* `volume_id` - The volume ID.

* `created_at` - The creation time of the volume.
