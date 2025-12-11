---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_instances"
description: |-
  Use this data source to get the list of ESW instances.
---

# huaweicloud_esw_instances

Use this data source to get the list of ESW instances.

## Example Usage

```hcl
data "huaweicloud_esw_instances" "test" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instances. If omitted, the provider-level region will
  be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance.

* `name` - (Optional, String) Specifies the name of the instance.

* `description` - (Optional, String) Specifies the description of the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of instances.
  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the name of the instance.

* `project_id` - Indicates the project ID.

* `region` - Indicates the region.

* `flavor_ref` - Indicates the flavor of the instance.

* `ha_mode` - Indicates the high availability mode of the instance.

* `status` - Indicates the status of the instance.

* `created_at` - Indicates the created time of the instance.

* `updated_at` - Indicates the updated time of the instance.

* `description` - Indicates the description of the instance.

* `availability_zones` - Indicates the availability zones of the instance.
  The [availability_zones](#availability_zones_struct) structure is documented below.

* `tunnel_info` - Indicates the local tunnel information of the instance.
  The [tunnel_info](#tunnel_info_struct) structure is documented below.

* `charge_infos` - Indicates the charge infos of the instance.
  The [charge_infos](#charge_infos_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `primary` - Indicates the availability zones where the default primary node is located.

* `standby` - Indicates the availability zones where the default standby node is located.

<a name="tunnel_info_struct"></a>
The `tunnel_info` block supports:

* `vpc_id` - Indicates the vpc ID.

* `virsubnet_id` - Indicates the subnet ID.

* `tunnel_ip` - Indicates the tunnel IP.

* `tunnel_port` - Indicates the tunnel port.

* `tunnel_type` - Indicates the tunnel type.

<a name="charge_infos_struct"></a>
The `charge_infos` block supports:

* `charge_mode` - Indicates the charge mode.
