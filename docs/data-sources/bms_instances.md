---
subcategory: "Bare Metal Server (BMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_bms_instances"
description: ""
---

# huaweicloud_bms_instances

Use this data source to get the list of BMS instances.

## Example Usage

```hcl
data "huaweicloud_bms_instances" "test" {
  status = "active"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of BMS instance.

* `flavor_id` - (Optional, String) Specifies the flavor ID of the desired flavor for the BMS instance.

* `name` - (Optional, String) Specifies the name of the BMS instance.

* `status` - (Optional, String) Specifies the status of the instance.
  Value options are as follows:
  + **ACTIVE**: running, stopping, deleting.
  + **BUILD**: creating.
  + **ERROR**: faulty.
  + **HARD_REBOOT**: forcibly restarting.
  + **REBOOT**: restarting.
  + **SHUTOFF**: stopped, starting, deleting, rebuilding, reinstalling OS, OS reinstallation failed, frozen.

* `tags` - (Optional, String) Specifies the BMS tags. The value can be: **__type_baremetal** or other custom tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - The list of BMS instances.
  The [servers](#attrblock--servers) structure is documented below.

<a name="attrblock--servers"></a>
The `servers` block supports:

* `id` - The ID of the BMS instance.

* `name` - The name of the BMS instance.

* `availability_zone` - The availability zone in which to create the BMS instance.

* `disk_config` - The disk configuration.

* `enterprise_project_id` - The enterprise project ID of BMS instance.

* `flavor_id` - The ID of the BMS flavor.

* `flavor_name` - The name of the BMS flavor.

* `vcpus` - The number of vCPUs.

* `memory` - The memory size in GB.

* `disk` - The system disk size(GB) in the BMS flavor. The value `0` indicates that the disk size is not limited.

* `image_id` - The image ID of the desired image for the BMS instance.

* `vpc_id` - The ID of vpc in which to create the BMS instance.

* `agency_name` - The IAM agency name which is created on IAM to provide temporary credentials for BMS
  to access cloud services.

* `image_name` - The image name of the desired image for the BMS instance.

* `image_type` - The image type of the desired image for the BMS instance.

* `key_pair` - The name of a key pair for logging in to the BMS using key pair authentication.

* `launched_at` - The start time of the BMS instance.

* `locked` - Whether the BMS instance is locked.

* `nics` - The list of one or more networks to attach to the BMS instance.
  The [nics](#attrblock--servers--nics) structure is documented below.

* `root_device_name` - The device name of the BMS system disk.

* `security_groups` - The list of one or more security group IDs to associate with the BMS instance.

* `status` - The status of the BMS instance.

* `user_data` - The user data to be injected during the BMS instance creation.

* `user_id` - The user ID.

* `vm_state` - The stable status of the BMS instance.

* `created_at` - The creation time of the BMS instance.

* `updated_at` - The latest update time of the BMS instance.

* `description` - The description of the BMS instance.

* `tags` - The key/value pairs to associate with the BMS instance.

* `volumes_attached` - The list of disks attached to the BMS instance.
  The [volumes_attached](#attrblock--servers--volumes_attached) structure is documented below.

<a name="attrblock--servers--nics"></a>
The `nics` block supports:

* `subnet_id` - The ID of subnet to attach to the BMS instance.

* `ip_address` - The fixed IPv4 address to be used on this network.

* `mac_address` - The MAC address of the nic.

* `port_id` - The port ID corresponding to the IP address.

<a name="attrblock--servers--volumes_attached"></a>
The `volumes_attached` block supports:

* `id` - The disk ID in UUID format.

* `boot_index` - Whether it is a boot disk. `0` specifies a boot disk, and `-1` specifies a non-boot disk.

* `delete_on_termination` - Whether to delete the disk when deleting the BMS.

* `device` - The device name of the disk.
