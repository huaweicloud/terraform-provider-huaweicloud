---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pool"
description: |-
  Manages a Workspace desktop pool resource within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pool

Manages a Workspace desktop pool resource within HuaweiCloud.

-> Before creating Workspace desktop, ensure that the Workspace service has been registered.

## Example Usage

```hcl
variable "desktop_pool_name" {}
variable "product_id" {}
variable "image_id" {}
variable "vpc_id" {}
variable "subnet_ids"  {
  type = list(string)
}
variable "security_group_ids"  {
  type = list(string)
}
variable "data_volume_sizes"  {
  type = list(number)
}
variable "authorized_object_list" {
  type = list(object({
    object_id   = string  
    object_type = string
    object_name = string
    user_group  = string
  }))
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = var.desktop_pool_name
  type                          = "DYNAMIC"
  size                          = 1
  product_id                    = var.product_id
  image_type                    = "gold"
  image_id                      = var.image_id
  vpc_id                        = var.vpc_id
  subnet_ids                    = var.subnet_ids
  disconnected_retention_period = 10

  root_volume {
    type = "SAS"
    size = 80
  }

  dynamic "security_groups" {
    for_each = var.security_group_ids
  
    content {
      id = security_groups.value
    }
  }

  dynamic "data_volumes" {
    for_each = var.data_volume_sizes
  
    content {
      type = "SAS"
      size = data_volumes.value
    }
  }

  dynamic "authorized_objects" {
    for_each = var.authorized_object_list
  
    content {
      object_id   = authorized_objects.value.object_id
      object_type = authorized_objects.value.object_type
      object_name = authorized_objects.value.object_name
      user_group  = authorized_objects.value.user_group
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the desktop pool.  
  The name valid length is limited from `1` to `15`, only Chinese and English characters, digits and hyphens (-) are allowed.

* `type` - (Required, String, NonUpdatable) Specifies the type of the desktop pool.  
  The valid values are as follows:
  + **DYNAMIC**
  + **STATIC**

* `size` - (Required, Int, NonUpdatable) Specifies the number of the desktops under the desktop pool.  
  The valid value ranges from `1` to `100`.

  -> This parameter will be affected by `autoscale_policy` parameter and will appear when `terraform plan` or
    `terraform apply`. If it is inconsistent with the script configuration, it can be ignored by `ignore_changes`
    in non-change scenarios.

* `product_id` - (Required, String, NonUpdatable) Specifies the specification ID of the desktop pool.

* `image_type` - (Required, String, NonUpdatable) Specifies the image type of the desktop pool.  
  The valid values are as follows:
  + **private**
  + **gold**

* `image_id` - (Required, String, NonUpdatable) Specifies the image ID of the desktop pool.

* `root_volume` - (Required, List, NonUpdatable) Specifies the system volume configuration of the desktop pool.  
  The [root_volume](#desktop_pool_volume) structure is documented below.

* `subnet_ids` - (Required, List, NonUpdatable) Specifies the list of the subnet IDs to which the desktop pool belongs.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the desktop pool belongs.

* `security_groups` - (Optional, List, NonUpdatable) Specifies the list of the security groups to which the
  desktop pool belongs.  
  The [security_groups](#desktop_pool_security_groups) structure is documented below.

* `data_volumes` - (Optional, List, NonUpdatable) Specifies the list of the data volume configurations of
  the desktop pool.  
  The [data_volumes](#desktop_pool_volume) structure is documented below.
  
* `authorized_objects` - (Optional, List, NonUpdatable) Specifies the list of the users or user groups
  to be authorized.  
  The [authorized_objects](#desktop_pool_authorized_objects) structure is documented below.

* `availability_zone` - (Optional, String) Specifies the availability zone to which the desktop pool belongs.

* `disconnected_retention_period` - (Optional, Int) Specifies the desktops and users disconnection retention period
  under desktop pool, in minutes.  
  The valid value ranges from `10` to `43,200`.  
  This parameter is available and required only when the `type` is set to **DYNAMIC**.

* `enable_autoscale` - (Optional, Bool) Specifies whether to enable elastic scaling of the desktop pool.  
  Defaults to **false**.

* `autoscale_policy` - (Optional, List) Specifies the automatic scaling policy of the desktop pool.  
  The [autoscale_policy](#desktop_pool_autoscale_policy) structure is documented below.  
  This parameter is available only when the `enable_autoscale` is set to **true**.

* `desktop_name_policy_id` - (Optional, String) Specifies the ID of the policy to generate the desktop name.

* `ou_name` - (Optional, String) Specifies the OU name corresponding to the AD server.  
  This parameter is available only when the AD server is connected.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which
  the desktop pool belongs.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

* `description` - (Optional, String) Specifies the description of the desktop pool.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the desktop pool.

* `in_maintenance_mode` - (Optional, Bool) Specifies whether to enable maintenance mode of the desktop pool.  
  Defaults to **false**.

<a name="desktop_pool_volume"></a>
The `root_volume` block supports:

* `type` - (Required, String) Specifies the type of the volume.  
  The valid values are as follows:
  + **SAS**: High I/O disk type.
  + **SSD**: Ultra-high I/O disk type.

* `size` - (Required, Int) Specifies the size of the volume, in GB.
  + For root volume, the valid value ranges from `80` to `1,020`.
  + For data volume, the valid value ranges from `10` to `8,200`.

<a name="desktop_pool_authorized_objects"></a>
The `authorized_objects` block supports:

* `object_id` - (Required, String) Specifies the ID of the object.

* `object_type` - (Required, String) Specifies the type of the object.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**

* `object_name` - (Required, String) Specifies the name of the object.

* `user_group` - (Required, String) Specifies the permission group to which the user belongs.  
  The valid values are as follows:
  + **sudo**: Linux administrator group.
  + **default**: Linux default user group.
  + **administrators**: Windows administrator group.
  + **users**: Windows standard user group.

<a name="desktop_pool_autoscale_policy"></a>
The `autoscale_policy` block supports:

* `autoscale_type` - (Optional, String) Specifies the type of automatic scaling policy.  
  The valid values are as follows:
  + **ACCESS_CREATED**: Create desktops during accessing.
  + **AUTO_CREATED**: Pre-creation desktops.

* `max_auto_created` - (Optional, Int) Specifies the maximum number of automatically created desktops.  
  The valid value ranges from `1` to `1,000`.

* `min_idle` - (Optional, Int) Specifies the desktops will be automatically created when the number of idle desktops is
  less than this value.  
  The valid value ranges from `1` to `1,000`.

<a name="desktop_pool_security_groups"></a>
The `security_groups` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the ID of the security group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also desktop pool ID.

* `root_volume` - The system volume configuration of the desktop pool.  
  The [root_volume](#attr_desktop_pool_volume) structure is documented below.
  
* `data_volumes` - The list of the data volume configurations of the desktop pool.  
  The [data_volumes](#attr_desktop_pool_volume) structure is documented below.

* `status` - The status of the desktop pool.
  + **STEADY**
  + **TEMPORARY**
  + **EXIST_FROZEN**
  + **UNKNOWN**
  
* `created_time` - The creation time of the desktop pool, in UTC format.

* `desktop_used` - The number of desktops associated with the users under the desktop pool.

* `product` - The product information of the desktop pool.  
  The [product](#attr_desktop_pool_product) structure is documented below.
  
* `image_name` - The image name of the desktop pool.

* `image_os_type` - The image OS type of the desktop pool.

* `image_os_version` - The image OS version of the desktop pool.

* `image_os_platform` - The image OS platform of the desktop pool.

<a name="attr_desktop_pool_volume"></a>
The `data_volumes` block supports:

* `id` - The ID of the volume.

<a name="attr_desktop_pool_product"></a>
The `product` block supports:

* `flavor_id` - The product specification ID of the desktop pool.

* `type` - The product type of the desktop pool.
  + **BASE**: Basic package of the product.

* `memory` - The product memory of the desktop pool.

* `cpu` - The product CPU of the desktop pool.

* `descriptions` - The product description of the desktop pool.

* `charging_mode` - The product charging mode of the desktop pool.
  + **0**: The yearly/monthly billing mode.
  + **1**: The pay-per-use billing mode.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The desktop pool can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_desktop_pool.test <id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `vpc_id`, `image_type`, `tags`, `ou_name`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_desktop_pool" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vpc_id, image_type, tags, ou_name,
    ]
  }
}
```
