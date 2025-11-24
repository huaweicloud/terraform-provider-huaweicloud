---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group"
description: |-
  Manages a Workspace APP server group resource within HuaweiCloud.
---

# huaweicloud_workspace_app_server_group

Manages a Workspace APP server group resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_group_name" {}
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "image_id" {}
variable "image_product_id" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = var.server_group_name
  os_type          = "Windows"
  flavor_id        = var.flavor_id
  vpc_id           = var.vpc_id
  subnet_id        = var.subnet_id
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = false
  image_id         = var.image_id
  image_type       = "gold"
  image_product_id = var.image_product_id
  app_type         = "SESSION_DESKTOP_APP"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the server group.  
  The name valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_) and
  hyphens (-) are allowed.

* `os_type` - (Required, String, ForceNew) Specifies the operating system type of the server group.
  Changing this creates a new resource.  
  Currently, only **Windows** is supported.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID of the server group.
  Changing this creates a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the server group belongs.
  Changing this creates a new resource.  
  This parameter value must be the VPC ID corresponding to the Workspace service.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID to which the server group belongs.
  Changing this creates a new resource.  
  This parameter value must be the subnet ID corresponding to the Workspace service.

* `system_disk_type` - (Required, String) Specifies the type of system disk.  
  The valid values are as follows:
  + **ESSD**: Extreme SSD type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SAS**: High I/O type.
  + **SATA**: Common I/O type.

* `system_disk_size` - (Required, Int) Specifies the size of system disk, in GB.  
  The minimum value of this parameter cannot be less than the disk size corresponding to the image.

* `image_id` - (Required, String) Specifies the image ID of the server group.

* `image_type` - (Required, String) Specifies the image type of the server group.  
  The valid values are as follows:
  + **gold**: The market image.
  + **public**: The public image.
  + **private**: The private image.
  + **shared**: The shared image.
  + **other**

* `image_product_id` - (Optional, String) Specifies the image product ID of the server group.  
  This parameter is required wnen the `image_type` parameter is set to **gold**.

* `is_vdi` - (Optional, Bool, ForceNew) Specifies the session mode of the server group. Defaults to **false**.
  + **false**: Multi-session mode.
  + **true**: Single-session mode.

  If the AD server is not connected, only the single-session mode is supported.  
  Changing this creates a new resource.

* `app_type` - (Optional, String) Specifies the type of application group associated with the server group.
  Defaults to **COMMON_APP**.  
  The valid values are as follows:
  + **SESSION_DESKTOP_APP**
  + **COMMON_APP**

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone of the server group.  
  If omitted, the AZ randomly assigned by the system is used.  
  Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the server group.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the server group.
  Supports up to `20` tags.
  
* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  server group belong.  
  This field is only valid for enterprise users, if omitted, default enterprise project will be used.  
  Changing this creates a new resource.

* `ip_virtual` - (Optional, List, ForceNew) Specifies the IP virtualization function configuration.
  Changing this creates a new resource.  
  This parameter is available only wnen the `is_vdi` parameter is set to **false**.  
  The [ip_virtual](#app_server_group_ip_virtual) structure is documented below.  

  -> If IP virtualization is enabled, each session is assigned a different IP address. The number of virtual IP addresses
    pre-allocated by the server is the same as the maximum number of sessions.

* `route_policy` - (Optional, List) Specifies the session scheduling policy of the server group.
  This parameter is available only wnen the `is_vdi` parameter is set to **false**.  
  The [route_policy](#app_server_group_route_policy) structure is documented below.

  -> If any metric of the server exceeds the threshold, new sessions will be rejected. The sessions will
     be automatically scheduled to other available servers.

* `ou_name` - (Optional, String) Specifies the OU name corresponding to the AD server.
  This parameter is available only when the AD server is connected.

* `extra_session_type` - (Optional, String, ForceNew) Specifies the additional session type.
  This parameter is available only wnen the `is_vdi` parameter is set to **false**.  
  Changing this creates a new resource.  
  The valid values are as follows:
  + **GPU**
  + **CPU**
  
* `extra_session_size` - (Optional, Int, ForceNew) Specifies the number of additional sessions for a single server.
  Changing this creates a new resource.  
  This parameter is available only wnen the `is_vdi` parameter is set to **false**.  
  The `extra_session_size` must be used together with `extra_session_type`.  
  The upper limit of the number of additional sessions for a single server is `10` times the number of vCPUs in the server
  specification minus the default number of sessions in the package.

* `primary_server_group_id` - (Optional, String, ForceNew) Specifies the ID of the primary server group.
  Changing this creates a new resource.

  -> 1. If this parameter is specified, the standby server is created.
     <br>2. The `os_type`, `is_vdi`, `app_type`, `ip_virtual` and `storage_mount_policy` parameters of the primary and
     standby server groups must be consistent.
     <br>3. After the `app_type` and `storage_mount_policy` parameters of the primary server group are changed, the change
     is automatically applied to the standby server group.

* `enabled` - (Optional, Bool) Whether to enable server group. Defaults to **true**.

* `storage_mount_policy` - (Optional, String) Specifies the NAS storage directory mounting policy on the APS.
  + **USER**: Only mount personal directories.
  + **SHARE**: Only mount shared directories.
  + **ANY**: No restrictions on the mounted directories (both personal and shared NAS storage directories will be
    automatically mounted).

<a name="app_server_group_ip_virtual"></a>
The `ip_virtual` block supports:

* `enable` - (Required, Bool, ForceNew) Whether to enable IP virtualization. Defaults to **false**.  
  Changing this creates a new resource.

<a name="app_server_group_route_policy"></a>
The `route_policy` block supports:

* `max_session` - (Optional, Int) The number of session connections of the server.  
  The maximum number of sessions is equal to the default number of sessions plus the number of additional sessions.

* `cpu_threshold` - (Optional, Int) The CPU usage of the server. The unit is `%`.  
  The valid value ranges from `1` to `100`.

* `mem_threshold` - (Optional, Int) The memory usage of the server. The unit is `%`.  
  The valid value ranges from `1` to `100`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also server group ID.

* `project_name` - The name of the project.

* `image_min_disk` - The minimum memory required to run the image, in MB. The default value is 0.

* `flavors` - The list of server flavors.

  The [flavors](#app_server_group_flavors) structure is documented below.

<a name="app_server_group_flavors"></a>
The `flavors` block supports:

* `id` - The ID of the flavor.

* `links` - The quick link information for relevant tags corresponding to server specifications.

  The [links](#app_server_group_flavor_links) structure is documented below.

<a name="app_server_group_flavor_links"></a>
The `links` block supports:

* `rel` - The shortcut link tag name.

* `href` - The corresponding shortcut link.

## Import

The server group resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_server_group.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `vpc_id`, `image_type`, `image_product_id`, `availability_zone`, `ip_virtual` and `route_policy`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_server_group" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vpc_id, image_type, image_product_id, availability_zone, ip_virtual, route_policy,
    ]
  }
}
```
