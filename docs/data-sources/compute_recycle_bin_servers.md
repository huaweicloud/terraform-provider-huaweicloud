---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_recycle_bin_servers"
description: |-
  Use this data source to get the list of servers in the recycle bin.
---

# huaweicloud_compute_recycle_bin_servers

Use this data source to get the list of servers in the recycle bin.

## Example Usage

```hcl
data "huaweicloud_compute_recycle_bin_servers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the ECS name, which is fuzzy-matched.

* `all_tenants` - (Optional, String) Specifies whether query ECSs for all tenants.
  Value options:
  + **1**: query VMs for all tenants
  + **0**: query VMs for current tenant

* `availability_zone` - (Optional, String) Specifies the availability zone.

* `expect_fields` - (Optional, String) Specifies the controls the query output.
  It controls whether to query based on the default field.
  Value options:
  + **launched_at**: the time when an ECS was started
  + **key_name**: the key pair that is used to authenticate an ECS
  + **locked**: whether an ECS is locked
  + **root_device_name**: the device name of the ECS system disk
  + **tenancy**: creating ECSs on a DeH or in a shared pool
  + **dedicated_host_id**: DeH ID
  + **enterprise_project_id**: querying the ECSs that are associated with an enterprise project
  + **tags**: list of ECS tags
  + **metadata**: ECS metadata
  + **addresses**: network addresses of an ECS
  + **security_groups**: information about the security group associated with the ECS
  + **volumes_attached**: information about the disks attached to the ECS
  + **image**: ECS image
  + **power_state**: ECS power status
  + **cpu_options**: custom CPU options
  + **market_info**: ECS billing information, including the billing type and expiration time

* `ip_address` - (Optional, String) Specifies the IP address.

* `tags` - (Optional, List) Specifies the queries ECSs with tags containing the specified value.

* `tags_key` - (Optional, List) Specifies the queries ECSs with tags containing the specified key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `servers` - Indicates the ECS information.

  The [servers](#servers_struct) structure is documented below.

<a name="servers_struct"></a>
The `servers` block supports:

* `os_scheduler_hints` - Indicates the ECS scheduling information.

  The [os_scheduler_hints](#servers_os_scheduler_hints_struct) structure is documented below.

* `sys_tags` - Indicates the ECS system tags.

  The [sys_tags](#servers_sys_tags_struct) structure is documented below.

* `os_ext_sts_task_state` - Indicates the ECS task status.

* `os_ext_sts_power_state` - Indicates the power status of the ECS.

* `os_srv_usg_terminated_at` - Indicates the time when the ECS was deleted.

* `description` - Indicates the ECS description.

* `progress` - Indicates the ECS progress.

* `metadata` - Indicates the ECS metadata

* `os_ext_srv_attr_root_device_name` - Indicates the device name of the ECS system disk.

* `locked` - Indicates whether an ECS is locked.

* `flavor` - Indicates the ECS flavor.

  The [flavor](#servers_flavor_struct) structure is documented below.

* `user_id` - Indicates the ID of the user for creating the ECS.

* `name` - Indicates the ECS name.

* `config_drive` - Indicates the config drive.

* `security_options` - Indicates the security options.

  The [security_options](#servers_security_options_struct) structure is documented below.

* `hypervisor` - Indicates the virtualization information.

  The [hypervisor](#servers_hypervisor_struct) structure is documented below.

* `host_id` - Indicates the ID of the host where the ECS located.

* `os_ext_srv_attr_hypervisor_hostname` - Indicates the name of the host on which the ECS is deployed.

* `id` - Indicates the ECS ID in UUID format.

* `os_ext_srv_attr_reservation_id` - Indicates the ID reserved for the ECSs to be created in a batch.

* `os_ext_srv_attr_launch_index` - Indicates the sequence in which ECSs start if the ECSs are created in a batch.

* `tags` - Indicates ECS tags.

* `os_ext_srv_attr_instance_name` - Indicates the ECS alias.

* `security_groups` - Indicates the security groups of the ECS.

  The [security_groups](#servers_security_groups_struct) structure is documented below.

* `os_extended_volumes_volumes_attached` - Indicates the disks attached to an ECS.

  The [os_extended_volumes_volumes_attached](#servers_os_extended_volumes_volumes_attached_struct) structure is
  documented below.

* `cpu_options` - Indicates the CPU options.

  The [cpu_options](#servers_cpu_options_struct) structure is documented below.

* `os_ext_az_availability_zone` - Indicates the AZ of an ECS.

* `host_status` - Indicates the status of the host accommodating the ECS.

* `os_ext_srv_attr_ramdisk_id` - Indicates the UUID of the ramdisk image if an AMI image is used.

* `os_ext_srv_attr_kernel_id` - Indicates the UUID of the kernel image if an AMI image is used.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `updated` - Indicates the last time when the ECS was updated.

* `auto_terminate_time` - Indicates the scheduled deletion time for the ECS.

* `os_ext_sts_vm_state` - Indicates the ECS status.

* `os_ext_srv_attr_hostname` - Indicates the host name of the ECS.

* `os_ext_srv_attr_host` - Indicates the name of the host where the ECS is deployed.

* `os_dcf_disk_config` - Indicates the disk configuration type.

* `os_srv_usg_launched_at` - Indicates the time when the ECS was started.

* `os_ext_srv_attr_user_data` - Indicates the user data (encoded) configured during ECS creation.

* `status` - Indicates  the ECS status.

* `key_name` - Indicates the key pair that is used to authenticate an ECS.

* `image` - Indicates the ECS image.

  The [image](#servers_image_struct) structure is documented below.

* `created` - Indicates the time when the ECS was created.

* `network_interfaces` - Indicates the network interface information.

  The [network_interfaces](#servers_network_interfaces_struct) structure is documented below.

<a name="servers_os_scheduler_hints_struct"></a>
The `os_scheduler_hints` block supports:

* `group` - Indicates the ECS group ID in UUID format.

* `tenancy` - Indicates creates ECSs on a dedicated or shared host.

* `dedicated_host_id` - Indicates the dedicated host ID.

<a name="servers_sys_tags_struct"></a>
The `sys_tags` block supports:

* `key` - Indicates the system tag key.

* `value` - Indicates the system tag value.

<a name="servers_flavor_struct"></a>
The `flavor` block supports:

* `ram` - Indicates the memory size (MiB) in the ECS flavor.

* `gpus` - Indicates the GPU information in the ECS flavor.

  The [gpus](#flavor_gpus_struct) structure is documented below.

* `asic_accelerators` - Indicates the ASIC information in the ECS flavor.

  The [asic_accelerators](#flavor_asic_accelerators_struct) structure is documented below.

* `id` - Indicates the ECS flavor ID.

* `name` - Indicates the ECS flavor name.

* `disk` - Indicates the system disk size in the ECS flavor.
  Value 0 indicates that the disk size is not limited.

* `vcpus` - Indicates the number of vCPUs in the ECS flavor.

<a name="flavor_gpus_struct"></a>
The `gpus` block supports:

* `count` - Indicates the number of GPUs.

* `memory_mb` - Indicates the GPU memory size, in MB.

* `name` - Indicates the GPU name.

<a name="flavor_asic_accelerators_struct"></a>
The `asic_accelerators` block supports:

* `memory_mb` - Indicates the ASIC memory size, in MB.

* `name` - Indicates the ASIC name.

* `count` - Indicates the number of ASICs.

<a name="servers_security_options_struct"></a>
The `security_options` block supports:

* `secure_boot_enabled` - Indicates whether support secure boot.

* `tpm_enabled` - Indicates whether support vtpm start.

<a name="servers_hypervisor_struct"></a>
The `hypervisor` block supports:

* `hypervisor_type` - Indicates the virtualization type.

* `csd_hypervisor` - Indicates the hypervisor csd info.

<a name="servers_security_groups_struct"></a>
The `security_groups` block supports:

* `name` - Indicates the security group name or UUID.

* `id` - Indicates the security group ID.

<a name="servers_os_extended_volumes_volumes_attached_struct"></a>
The `os_extended_volumes_volumes_attached` block supports:

* `id` - Indicates the disk ID in UUID format.

* `delete_on_termination` - Indicates whether the disk is deleted with the ECS.

* `boot_index` - Indicates the EVS disk boot sequence.

* `device` - Indicates the drive letter of the EVS disk.
  Which is the device name of the EVS disk.

<a name="servers_cpu_options_struct"></a>
The `cpu_options` block supports:

* `hw_cpu_threads` - Indicates whether to enable CPU hyper-threading.

<a name="servers_image_struct"></a>
The `image` block supports:

* `id` - Indicates the image ID.

<a name="servers_network_interfaces_struct"></a>
The `network_interfaces` block supports:

* `port_id` - Indicates the port ID.

* `primary` - Indicates whether the network interface is a primary network interface.

* `ip_addresses` - Indicates private IPv4 addresses.

* `ipv6_addresses` - Indicates private IPv6 addresses.

* `association` - Indicates information about the associated EIP.

  The [association](#network_interfaces_association_struct) structure is documented below.

* `subnet_id` - Indicates the subnet ID.

<a name="network_interfaces_association_struct"></a>
The `association` block supports:

* `public_ip_address` - Indicates the IPv4 address of the EIP.
