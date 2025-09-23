---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_image_server"
description: |-
  Manages an image server resource of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_image_server

Manages an image server resource of Workspace APP within HuaweiCloud.

## Example Usage

```hcl
variable "image_server_name" {}
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}
variable "image_id" {}
variable "image_spec_code" {}
variable "image_source_product_id" {}
variable "user_name" {}

resource "huaweicloud_workspace_app_image_server" "test" {
  name                    = var.image_server_name
  flavor_id               = var.flavor_id
  vpc_id                  = var.vpc_id
  subnet_id               = var.subnet_id
  image_id                = var.image_id
  image_type              = "gold"
  spec_code               = var.image_spec_code
  image_source_product_id = var.image_source_product_id
  is_vdi                  = true

  authorize_accounts {
    account = var.user_name
    type    = "USER"
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  is_delete_associated_resources = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the image server.
  Changing this creates a new resource.  
  The name valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_) and
  hyphens (-) are allowed and cannot contain spaces.
  
* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID of the image server.
  Changing this creates a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the image server belongs.
  Changing this creates a new resource.  
  This parameter value must be the VPC ID corresponding to the Workspace service.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID to which the image server belongs.
  Changing this creates a new resource.  
  This parameter value must be the VPC ID corresponding to the Workspace service.

* `root_volume` - (Required, List, ForceNew) Specifies the system disk configuration of the image server.
  Changing this creates a new resource.
  The [root_volume](#app_image_server_root_volume) structure is documented below.
  
* `authorize_accounts` - (Required, List, ForceNew) Specifies the list of the management accounts for creating the image.
  Changing this creates a new resource.
  The [authorize_accounts](#app_image_server_authorize_accounts) structure is documented below.

* `image_id` - (Required, String, ForceNew) Specifies the basic image ID of the image server.
  Changing this creates a new resource.

* `image_type` - (Required, String, ForceNew) Specifies the basic image type of the image server.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **gold**: The market image.
  + **public**: The public image.
  + **private**: The private image.
  + **shared**: The shared image.
  + **other**

* `spec_code` - (Optional, String, ForceNew) Specifies the specification code of the basic image to which the image
  server belongs. Changing this creates a new resource.  
  This parameter is required when the `image_type` parameter is set to **gold**.

* `image_source_product_id` - (Optional, String, ForceNew) Specifies the basic image product ID of the image server.
  Changing this creates a new resource.  
  This parameter is required when the `image_type` parameter is set to **gold**.
  
* `is_vdi` - (Optional, Bool, ForceNew) Specifies the session mode of the image server.
  Changing this creates a new resource.
  + **false**: Multi-session mode (default value).
  + **true**: Single-session mode.

  If the AD server is not connected, only the single-session mode is supported.  
  
* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone of the image server.
  Changing this creates a new resource.  
  If omitted, the AZ randomly assigned by the system is used.  

* `description` - (Optional, String) Specifies the description of the image server.

* `ou_name` - (Optional, String, ForceNew) Specifies the OU name corresponding to the AD server.
  Changing this creates a new resource.  
  This parameter is available only when the AD server is connected.

* `extra_session_type` - (Optional, String, ForceNew) Specifies the additional session type.
  Changing this creates a new resource.  
  This parameter is available only when the `is_vdi` parameter is set to **false**.  
  The valid values are as follows:
  + **GPU**
  + **CPU**

* `extra_session_size` - (Optional, Int, ForceNew) Specifies the number of additional sessions for a single server.
  Changing this creates a new resource.  
  This parameter is available only when the `is_vdi` parameter is set to **false**.  
  The `extra_session_size` must be used together with `extra_session_type`.  
  The upper limit of the number of additional sessions for a single server is `10` times the number of vCPUs in the server
  specification minus the default number of sessions in the package.

* `route_policy` - (Optional, List, ForceNew) Specifies the session scheduling policy of the server associated with
  the image server. Changing this creates a new resource.  
  This parameter is available only wnen the `is_vdi` parameter is set to **false**.  
  The [route_policy](#app_image_server_route_policy) structure is documented below.

  -> If any metric of the server exceeds the threshold, new sessions will be rejected. The sessions will
     be automatically scheduled to other available servers.

* `scheduler_hints` - (Optional, List, ForceNew) Specifies the configuration of the dedicate host.
  Changing this creates a new resource.
  The [scheduler_hints](#app_image_server_scheduler_hints) structure is documented below.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the image server.
  Supports up to `20` tags.  
  Changing this creates a new resource.
  
* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the image
  server belong. Changing this creates a new resource.  
  This parameter is only valid for enterprise users, if omitted, default enterprise project will be used.

* `is_delete_associated_resources` - (Optional, Bool) Specifies whether to delete resources associated with this image server
  after deleting it, defaults to **false**.

  -> If this parameter is set to **true**, deleting the resource will also delete the associated server group, server
     and application group resources, but the image product related resources will be retained.
  
<a name="app_image_server_authorize_accounts"></a>
The `authorize_accounts` block supports:

* `account` - (Required, String, ForceNew) Specifies the name of the account.
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the account.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **USER**

* `domain` - (Optional, String, ForceNew) Specifies the domain name of the Workspace service.
  Changing this creates a new resource.

<a name="app_image_server_root_volume"></a>
The `root_volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the disk type of the image server.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **ESSD**: Extreme SSD type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SAS**: High I/O type.
  + **SATA**: Common I/O type.

* `size` - (Required, Int, ForceNew) Specifies the disk size of the image server, in GB.
  Changing this creates a new resource.  
  The system disk size must be sufficient for the basic image and the application to be installed.

<a name="app_image_server_route_policy"></a>
The `route_policy` block supports:

* `max_session` - (Optional, Int, ForceNew) Specifies the number of session connections of the server.
  Changing this creates a new resource.  
  The maximum number of sessions is equal to the default number of sessions plus the number of additional sessions.

* `cpu_threshold` - (Optional, Int, ForceNew) Specifies the CPU usage of the server. The unit is `%`.  
  Changing this creates a new resource.  
  The valid value ranges from `1` to `100`.

* `mem_threshold` - (Optional, Int, ForceNew) Specifies the memory usage of the server. The unit is `%`.  
  Changing this creates a new resource.  
  The valid value ranges from `1` to `100`.

<a name="app_image_server_scheduler_hints"></a>
The `scheduler_hints` block supports:

* `dedicated_host_id` - (Optional, String, ForceNew) Specifies the ID of the dedicate host.
  Changing this creates a new resource.

* `tenancy` - (Optional, String, ForceNew) Specifies the type of the dedicate host.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also iamge server ID.

* `created_at` - The creation time of the image server, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 20 minutes.

## Import

The image server resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_image_server.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `flavor_id`, `vpc_id`, `subnet_id`, `root_volume`, `image_source_product_id`, `is_vdi`,
`availability_zone`, `ou_name`, `extra_session_type`, `extra_session_size`, `route_policy`, `scheduler_hints`, `tags`,
`enterprise_project_id`,  `is_delete_associated_resources`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_image_server" "test" {
  ...

  lifecycle {
    ignore_changes = [
      flavor_id, vpc_id, subnet_id, root_volume, image_source_product_id, is_vdi, availability_zone, ou_name, extra_session_type,
      extra_session_size, route_policy, scheduler_hints, tags, enterprise_project_id, is_delete_associated_resources,
    ]
  }
}
```
