---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_instances"
description: ""
---

# huaweicloud_compute_instances

Use this data source to get the list of the compute instances.

## Example Usage

```hcl
variable "name_regex" {}

data "huaweicloud_compute_instances" "test" {
  name = var.name_regex
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the instances.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the instance name, which can be queried with a regular expression.
  The instance name supports fuzzy matching query too.

* `instance_id` - (Optional, String) Specifies the ECS ID.

* `flavor_name` - (Optional, String) Specifies the flavor name of the instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `status` - (Optional, String) Specifies the status of the instance. The valid values are as follows:
  + **ACTIVE**: The instance is running properly.
  + **SHUTOFF**: The instance has been properly stopped.
  + **ERROR**: An error has occurred on the instance.

* `image_id` - (Optional, String) Specifies the image ID of the instance.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `fixed_ip_v4` - (Optional, String)  Specifies the IPv4 addresses of the ECS.

* `availability_zone` - (Optional, String) Specifies the availability zone where the instance is located.
  Please following [reference](https://developer.huaweicloud.com/intl/en-us/endpoint?ECS) for this argument.

* `key_pair` - (Optional, String) Specifies the key pair that is used to authenticate the instance.

* `tags` - (Optional, Map) Specifies the tags to qurey the instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - List of ECS instance details. The object structure of each ECS instance is documented below.

The `instances` block supports:

* `id` - The instance ID in UUID format.

* `name` - The instance name.

* `image_id` - The image ID of the instance.

* `image_name` - The image name of the instance.

* `flavor_id` - The flavor ID.

* `flavor_name` - The flavor name of the instance.

* `enterprise_project_id` - The enterprise project ID.

* `status` - The instance status.

* `availability_zone` - The availability zone where the instance is located.

* `public_ip` - The EIP address that is associated to the instance.

* `system_disk_id` - The system disk volume ID.

* `key_pair` - The key pair that is used to authenticate the instance.

* `security_group_ids` - An array of one or more security group IDs to associate with the instance.

* `user_data` - The user data (information after encoding) configured during instance creation.

* `network` - An array of one or more networks to attach to the instance.
  The [network object](#compute_instances_network_object) structure is documented below.

* `volume_attached` - An array of one or more disks to attach to the instance.
  The [volume attached object](#compute_instances_volume_object) structure is documented below.

* `scheduler_hints` - The scheduler with hints on how the instance should be launched.
  The [scheduler hints](#compute_instances_scheduler_hint_object) structure is documented below.

* `tags` - The key/value pairs to associate with the instance.

* `charging_mode` - The charging mode of the instance. Valid values are **prePaid**, **postPaid** and **spot**.

* `expired_time` - The expired time of prePaid instance, in UTC format.

<a name="compute_instances_network_object"></a>
The `network` block supports:

* `uuid` - The network ID to attach to the server.
* `port` - The port ID corresponding to the IP address on that network.
* `mac` - The MAC address of the NIC on that network.
* `fixed_ip_v4` - The fixed IPv4 address of the instance on this network.
* `fixed_ip_v6` - The Fixed IPv6 address of the instance on that network.

<a name="compute_instances_volume_object"></a>
The `volume_attached` block supports:

* `volume_id` - The volume ID on that attachment.
* `boot_index` - The volume boot index on that attachment.
* `is_sys_volume` - Whether the volume is the system disk.
* `size` - The volume size on that attachment.
* `type` - The volume type on that attachment.
* `pci_address` - The volume pci address on that attachment.

<a name="compute_instances_scheduler_hint_object"></a>
The `scheduler_hints` block supports:

* `group` - The server group ID where the instance will be placed into.
