---
subcategory: "Workspace"
---

# huaweicloud_workspace_desktop

Manages a Workspace desktop resource within HuaweiCloud.

## Example Usage

### Create a desktop using market image

```hcl
variable "flavor_id" {}
variable "image_id" {}
variable "vpc_id" {}
variable "network_id" {}
variable "security_group_id" {}
variable "desktop_name" {}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroups" "test" {
  // Security group automatically created when first opening the Workspace account, do not remove
  name = "WorkspaceUserSecurityGroup"
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id  = var.flavor_id
  image_type = "market"
  image_id   = var.image_id

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = var.vpc_id
  security_groups   = setunion(data.huaweicloud_networking_secgroups.test.security_groups[*].id,
    [var.security_group_id])

  nics {
    network_id = var.network_id
  }

  name       = var.desktop_name
  user_name  = "TestUser"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Workspace desktop resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `flavor_id` - (Required, String) Specifies the flavor ID of desktop.

* `image_type` - (Required, String, ForceNew) Specifies the image type. The valid values are as follows:
  + **market**: The market image.
  + **gold**: The public image.
  + **private**: The private image.

  Changing this will create a new resource.

* `image_id` - (Required, String, ForceNew) Specifies the image ID to create the desktop.
  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the desktop belongs.
  Changing this will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the user name to which the desktop belongs.
  The name can contain `1` to `20` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter. Changing this will create a new resource.

* `user_email` - (Required, String, ForceNew) Specifies the user email.
  Some operations on the desktop (such as creation, deletion) will notify the user by sending an email.
  Changing this will create a new resource.

* `root_volume` - (Required, List) Specifies the configuration of system volume.
  The [object](#desktop_volume) structure is documented below.

* `data_volume` - (Optional, List) Specifies the configuration of data volumes.
  The [object](#desktop_volume) structure is documented below.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone where the desktop is located.
  Changing this will create a new resource.

* `security_groups` - (Optional, List, ForceNew) Specifies the ID list of security groups.
  In addition to the custom security group, it must also contain a security group called **WorkspaceUserSecurityGroup**.
  Changing this will create a new resource.

* `user_group` - (Required, String, ForceNew) Specifies the user group to which the desktop belongs.
  The valid values are as follows:
  + **sudo**: Linux administrator group.
  + **default**: Linux default user group.
  + **administrators**: Windows administrator group.
  + **users**: Windows standard user group.

  Changing this will create a new resource.

* `nic` - (Optional, List, ForceNew) Specifies the NIC information corresponding to the desktop.
  The [object](#desktop_nic) structure is documented below. Changing this will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the desktop name.
  The name can contain `1` to `15` characters, only letters, digits and hyphens (-) are allowed.
  The name must start with a letter or digit and cannot end with a hyphen.
  Changing this will create a new resource.

* `email_notification` - (Optional, Bool, ForceNew) Specifies whether to send emails to user mailbox during important
  operations. Changing this will create a new resource.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs of the desktop.
  Changing this will create a new resource.

* `delete_user` - (Optional, Bool) Specifies whether to delete user associated with this desktop after deleting it.
  The user can only be successfully deleted if the user has no other desktops.

<a name="desktop_volume"></a>
The `root_volume` and `data_volume` block supports:

* `type` - (Required, String) Specifies the type of system volume.
  The valid values are as follows:
  + **SAS**: High I/O disk type.
  + **SSD**: Ultra-high I/O disk type.

  -> Updates are not supported for this parameter. Changing this will not create a new resource, but will throw an
     error.

* `size` - (Required, Int) Specifies the size of system volume, in GB.
  + For root volume, the valid value is range from `80` to `1,020`.
  + For data volume, the valid value is range from `10` to `8,200`.

<a name="desktop_nic"></a>
The `nic` block supports:

* `network_id` - (Required, String, ForceNew) Specifies the network ID of subnet resource.
  Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The desktop ID in UUID format.

* `root_volume` - The configuration of system volume.
  The [object](#desktop_volume_attr) structure is documented below.

* `data_volume` - The configuration of data volumes.
  The [object](#desktop_volume_attr) structure is documented below.

<a name="desktop_volume_attr"></a>
The `root_volume` and `data_volume` block supports:

* `id` - The volume ID.

* `name` - The volume name.

* `device` - The device location to which the volume is attached.

* `created_at` - The time that the volume was created.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

Desktops can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_workspace_desktop.test 339d2539-e945-4090-a08d-c16badc0c6bb
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `nic` and `user_email`.
It is generally recommended running `terraform plan` after importing a desktop.
You can then decide if changes should be applied to the desktop, or the resource definition should be updated to
align with the desktop. Also you can ignore changes as below.

```
resource "huaweicloud_workspace_desktop" "test" {
  ...

  lifecycle {
    ignore_changes = [
      user_email, nic,
    ]
  }
}
```
