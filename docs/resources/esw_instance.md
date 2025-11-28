---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_instance"
description: |-
  Manages an ESW instance resource within HuaweiCloud.
---

# huaweicloud_esw_instance

Manages an ESW instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_id" {}
variable "virsubnet_id" {}

resource "huaweicloud_esw_instance" "test" {
  name        = "test_name"
  flavor_ref  = "l2cg.small.ha"
  ha_mode     = "ha"
  description = "terraform test description"

  availability_zones {
    primary = "cn-north-4a"
    standby = "cn-north-4b"
  }

  tunnel_info {
    vpc_id       = var.vpc_id
    virsubnet_id = var.virsubnet_id
    tunnel_ip    = "192.168.0.192"
  }

  charge_infos {
    charge_mode = "postPaid"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ESW instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the instance.

* `flavor_ref` - (Required, String, NonUpdatable) Specifies the flavor of the instance.

* `ha_mode` - (Required, String, NonUpdatable) Specifies the high availability mode of the instance. Value options: **ha**.

* `availability_zones` - (Required, List, NonUpdatable) Specifies the availability zones of the instance.
  The [availability_zones](#availability_zones_struct) structure is documented below.

* `tunnel_info` - (Required, List, NonUpdatable) Specifies the local tunnel information of the instance.
  The [tunnel_info](#tunnel_info_struct) structure is documented below.

* `charge_infos` - (Required, List, NonUpdatable) Specifies the charge infos of the instance.
  The [charge_infos](#charge_infos_struct) structure is documented below.

* `description` - (Optional, String) Specifies the description of the instance.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `primary` - (Required, String, NonUpdatable) Specifies the availability zones where the default primary node is located.

* `standby` - (Required, String, NonUpdatable) Specifies the availability zones where the default standby node is located.

<a name="tunnel_info_struct"></a>
The `tunnel_info` block supports:

* `vpc_id` - (Required, String, NonUpdatable) Specifies the vpc ID.

* `virsubnet_id` - (Required, String, NonUpdatable) Specifies the subnet ID.

* `tunnel_ip` - (Optional, String, NonUpdatable) Specifies the tunnel IP.

<a name="charge_infos_struct"></a>
The `charge_infos` block supports:

* `charge_mode` - (Required, String, NonUpdatable) Specifies the charge mode. Value options: **postPaid**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tunnel_info` - Indicates the local tunnel information of the instance.
  The [tunnel_info](#tunnel_info_attribute) structure is documented below.

* `status` - Indicates the status of the instance.

* `created_at` - Indicates the created time of the instance.

* `updated_at` - Indicates the updated time of the instance.

<a name="tunnel_info_attribute"></a>
The `tunnel_info` block supports:

* `tunnel_port` - Indicates the tunnel port.

* `tunnel_type` - Indicates the tunnel type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

This resource can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_esw_instance.test <id>
```
