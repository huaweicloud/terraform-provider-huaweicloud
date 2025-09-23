---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server"
description: |-
  Manages a Workspace APP server resource under specified server group within HuaweiCloud.
---

# huaweicloud_workspace_app_server

Manages a Workspace APP server resource under specified server group within HuaweiCloud.

## Example Usage

```hcl
variable "server_group_id" {}
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_workspace_app_server" "test" {
  server_group_id = var.server_group_id
  type            = "createApps"
  flavor_id       = var.flavor_id

  root_volume {
    type = "SAS"
    size = 80
  }

  vpc_id              = var.vpc_id
  subnet_id           = var.subnet_id
  update_access_agent = true

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `server_group_id` - (Required, String, ForceNew) Specifies the server group ID to which the server belongs.
  Changing this creates a new resource.

* `name` - (Optional, String) Specifies the name of the server.  
  The name valid length is limited from `1` to `64`, only Chinese and English characters, digits, underscores (_) and
  hyphens (-) are allowed.

* `type` - (Required, String, ForceNew) Specifies the type of the server.
  Changing this creates a new resource.  
  Currently, only **createApps** is supported.
  
* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID of the server.
  Changing this creates a new resource.  
  This parameter value must be consistent with the server group to which it belongs.

* `root_volume` - (Required, List) Specifies the system disk configuration of the server.  
  This parameter value must be consistent with the server group to which it belongs.  
  The [root_volume](#app_server_root_volume) structure is documented below.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the server belongs.
  Changing this creates a new resource.  
  This parameter value must be consistent with the server group to which it belongs.
  
* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID to which the server belongs.
  Changing this creates a new resource.  
  This parameter value must be consistent with the server group to which it belongs.

* `os_type` - (Optional, String, ForceNew) Specifies the operating system type of the server.
  Changing this creates a new resource.  
  Currently, only **Windows** is supported.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone of the server.
  Changing this creates a new resource.  
  If omitted, the AZ randomly assigned by the system is used.

* `description` - (Optional, String) Specifies the description of the server.

* `maintain_status` - (Optional, Bool) Specifies whether to enable maintenance mode. Defaults to **false**.

* `ou_name` - (Optional, String) Specifies the OU name corresponding to the AD server.  
  This parameter is available only when the AD server is connected.

* `update_access_agent` - (Optional, Bool) Specifies whether to automatically upgrade protocol component. Defaults to **false**.

* `scheduler_hints` - (Optional, List, ForceNew) Specifies the configuration of the dedicate host.
  Changing this creates a new resource.  
  This parameter is available only when `charging_mode` is set to **postPaid**.  
  The [scheduler_hints](#app_server_scheduler_hints) structure is documented below.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the server. Defaults to **postPaid**.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **prePaid**: the yearly/monthly billing mode.
  + **postPaid**: the pay-per-use billing mode.
  
* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the server.
  Changing this creates a new resource.  
  This parameter is required and available if `charging_mode` is set to **prePaid**.  
  The valid values are as follows:
  + **month**
  + **year**

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the server.
  Changing this creates a new resource.
  + If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  + If `period_unit` is set to **year**, the value ranges from `1` to `3`.

  This parameter is required and available if `charging_mode` is set to **prePaid**.  
  
* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Defaults to **false**.  
  This parameter is required and available if `charging_mode` is set to **prePaid**.  
  The valid values are **true** and **false**.
  
<a name="app_server_root_volume"></a>
The `root_volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the disk type of the server.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **ESSD**: Extreme SSD type.
  + **SSD**: Ultra-high I/O type.
  + **GPSSD**: General purpose SSD type.
  + **SAS**: High I/O type.
  + **SATA**: Common I/O type.

* `size` - (Required, Int, ForceNew) Specifies the disk size of the server, in GB.
  Changing this creates a new resource.

<a name="app_server_scheduler_hints"></a>
The `scheduler_hints` block supports:

* `dedicated_host_id` - (Optional, String, ForceNew) Specifies the ID of the dedicate host.
  Changing this creates a new resource.

* `tenancy` - (Optional, String, ForceNew) Specifies the type of the dedicate host.
  Changing this creates a new resource.  
  Currently, only **dedicated** is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also server ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 40 minutes.
* `delete` - Default is 10 minutes.

## Import

The APP server resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_app_server.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `type`, `vpc_id`, `subnet_id`, `update_access_agent`, `scheduler_hints`, `period_unit`,
`period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_app_server" "test" {
  ...

  lifecycle {
    ignore_changes = [
      type, vpc_id, subnet_id, update_access_agent, scheduler_hints, period_unit, period, auto_renew,
    ]
  }
}
```
