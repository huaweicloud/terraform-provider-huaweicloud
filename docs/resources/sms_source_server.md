---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_server"
description: |-
  Manages an SMS source server resource within HuaweiCloud.
---

# huaweicloud_sms_source_server

Manages an SMS source server resource within HuaweiCloud.

## Example Usage

```hcl
variable "ip" {}
variable "name" {}
variable "agent_version" {}
variable "netmask" {}
variable "gateway" {}
variable "mac" {}

resource "huaweicloud_sms_source_server" "test" {
  ip                 = var.ip
  name               = var.name
  os_type            = "LINUX"
  os_version         = "Ubuntu 18.04 server 64bit"
  firmware           = "BIOS"
  boot_loader        = "GRUB"
  has_rsync          = true
  paravirtualization = true
  cpu_quantity       = 2
  memory             = 4018196480
  agent_version      = var.agent_version

  disks {
    name            = "Disk 0"
    partition_style = "MBR"
    device_use      = "BOOT"
    size            = 85897247744
    used_size       = 42943137792
    physical_volumes {
      device_use  = "OS"
      file_system = "ext4"
      mount_point = "/"
      name        = "/dev/vda1"
      size        = 42943137792
      used_size   = 1071640576
    }
  }

  networks {
    name    = "eth0"
    ip      = var.ip
    netmask = var.netmask
    gateway = var.gateway
    mtu     = 1
    mac     = var.mac
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Optional, String, NonUpdatable) Specifies the IP address of the source server.

  -> This parameter is mandatory for registering the source server with SMS and optional for updating the information
  about the source server.

* `name` - (Optional, String) Specifies the source server name in SMS.

* `hostname` - (Optional, String, NonUpdatable) Specifies the hostname of the source server.

  -> This parameter is mandatory for registering the source server with SMS and optional for updating the information
  about the source server.

* `os_type` - (Optional, String, NonUpdatable) Specifies the OS type of the source server, which can be Windows or Linux.
  Values can be **Windows** and **Linux**.

  -> This parameter is mandatory for registering the source server with SMS and optional for updating the information
  about the source server.

* `os_version` - (Optional, String, NonUpdatable) Specifies the OS version.

  -> This parameter is mandatory for registering the source server with SMS and optional for updating the information
  about the source server.

* `virtualization_type` - (Optional, String, NonUpdatable) Specifies the OS virtualization type.

* `linux_block_check` - (Optional, String, NonUpdatable) Specifies the Linux block-level check.

* `firmware` - (Optional, String, NonUpdatable) Specifies the boot mode.
  Values can be **BIOS** and **UEFI**.

* `cpu_quantity` - (Optional, Int, NonUpdatable) Specifies the number of CPUs, the unit is vCPUs.

* `memory` - (Optional, Int, NonUpdatable) Specifies the memory size, the unit is MB.

* `disks` - (Optional, List) Specifies the disk information of the source server.
  The [disks](#disks_struct) structure is documented below.

* `btrfs_list` - (Optional, List) Specifies the information about Btrfs file systems on the source server.
  The [btrfs_list](#btrfs_list_struct) structure is documented below.

  -> This parameter is mandatory for Linux.

* `networks` - (Optional, List) Specifies the information about NIC on the source server.
  The [networks](#networks_struct) structure is documented below.

* `domain_id` - (Optional, String, NonUpdatable) Specifies the tenant domain ID.

* `has_rsync` - (Optional, Bool, NonUpdatable) Specifies whether rsync is installed.

  -> This parameter is mandatory for Linux.
* `paravirtualization` - (Optional, Bool, NonUpdatable) Specifies whether the source server is paravirtualized.

  -> This parameter is mandatory for Linux.

* `raw_devices` - (Optional, String, NonUpdatable) Specifies the list of raw devices.

  -> This parameter is mandatory for Linux.

* `driver_files` - (Optional, Bool, NonUpdatable) Specifies whether any driver files are missing.

  -> This parameter is mandatory for Windows.

* `system_services` - (Optional, Bool, NonUpdatable) Specifies whether there are abnormal services.

  -> This parameter is mandatory for Windows.

* `account_rights` - (Optional, Bool, NonUpdatable) Specifies whether the account has the required permissions.

  -> This parameter is mandatory for Windows.

* `boot_loader` - (Optional, String, NonUpdatable) Specifies the system boot loader.
  Values can be **GRUB** and **LILO**.

  -> This parameter is mandatory for Linux.

* `system_dir` - (Optional, String, NonUpdatable) Specifies the system directory.

  -> This parameter is mandatory for Windows.

* `volume_groups` - (Optional, List) Specifies the volume groups.
  The [volume_groups](#volume_groups_struct) structure is documented below.

  -> This parameter is mandatory for Linux.

* `agent_version` - (Optional, String, NonUpdatable) Specifies the agent version.

* `kernel_version` - (Optional, String, NonUpdatable) Specifies the kernel version.

* `migration_cycle` - (Optional, String) Specifies the current migration stage of the source server.
  Values can be as follows:
  + **cutovering**: The target server for the source server is being launched.
  + **cutovered**: The target server for the source server is launched.
  + **checking**: The check is in progress.
  + **setting**: The configuration is in progress.
  + **replicating**: The data is being replicated.
  + **syncing**: The incremental data is being synchronized.

* `state` - (Optional, String) Specifies the source server status.
  Values can be as follows:
  + **unavailable**: The source server fails the environment check.
  + **waiting**: The source server is waiting for migration.
  + **initialize**: The migration of the source server is being initialized.
  + **replicate**: The source server is being replicated.
  + **syncing**: The source server is being synchronized.
  + **stopping**: The migration of the source server is being stopped.
  + **stopped**: The migration of the source server is stopped.
  + **deleting**: The source server record is being deleted.
  + **error**: An error occurs during the migration of the source server.
  + **cloning**: The target server for the source server is being cloned.
  + **cutovering**: The target server for the source server is being launched.
  + **finished**: The target server for the source server is launched.
  + **clearing**: The snapshot resources are being cleared.
  + **cleared**: The snapshot resources have been cleared.
  + **clearfailed**: The snapshot resources fail to be cleared.
  + **premigready**: The migration drill is ready.
  + **premiging**: The migration drill is in progress.
  + **premiged**: The migration drill has been completed.
  + **premigfailed**: The migration drill fails.

* `oem_system` - (Optional, Bool, NonUpdatable) Specifies whether the OS is an OEM version (Windows).

* `start_type` - (Optional, String, NonUpdatable) Specifies the startup mode.
  Values can be **MANUAL**, **MGC** or an empty string ("").

* `io_read_wait` - (Optional, Int, NonUpdatable) Specifies the disk read latency, the unit is ms.

* `has_tc` - (Optional, Bool, NonUpdatable) Specifies whether TC is installed.

  -> This parameter is mandatory for Linux.

* `platform` - (Optional, String, NonUpdatable) Specifies the platform.
  Values can be as follows:
  + **hw**: Huawei Cloud.
  + **ali**: Alibaba Cloud.
  + **aws**: AWS.
  + **azure**: Microsoft Azure.
  + **gcp**: Google Cloud.
  + **tencent**: Tencent Cloud.
  + **vmware**: VMware.
  + **hyperv**: Hyper-V.
  + **other**: other providers.

* `migprojectid` - (Optional, String) Specifies the ID of the migration project to which the source server belongs after
  the modification.

* `copystate` - (Optional, String) Specifies the source server status.
  Values can be as follows:
  + **UNAVAILABLE**: The source server fails the environment check.
  + **WAITING**: The source server is waiting for migration.
  + **INIT**: The migration of the source server is being initialized.
  + **REPLICATE**: The source server is being replicated.
  + **SYNCING**: The source server is being synchronized.
  + **SYNCING**: The migration of the source server is being stopped.
  + **STOPPED**: The migration of the source server is stopped.
  + **DELETING**: The source server record is being deleted.
  + **ERROR**: An error occurs during the migration of the source server.
  + **CLONING**: The target server for the source server is being cloned.
  + **CUTOVERING**: The target server for the source server is being launched.
  + **FINISHED**: The target server for the source server is launched.
  + **CLEARING**: The snapshot resources are being cleared.
  + **CLEARED**: The snapshot resources have been cleared.
  + **CLEARFAILED**: The snapshot resources fail to be cleared.
  + **premigready**: The migration drill is ready.
  + **premiging**: The migration drill is in progress.
  + **premiged**: The migration drill has been completed.
  + **premigfailed**: The migration drill fails.

<a name="disks_struct"></a>
The `disks` block supports:

* `name` - (Required, String) Specifies the disk name.

* `device_use` - (Required, String) Specifies the disk function.

* `size` - (Required, Int) Specifies the disk size, the unit is bytes.

* `used_size` - (Required, Int) Specifies the used disk space, the unit is bytes.

* `physical_volumes` - (Required, List) Specifies the information about physical partitions on the disk.
  The [physical_volumes](#disks_physical_volumes_struct) structure is documented below.

* `partition_style` - (Optional, String) Specifies the disk partition type.
  Values can be as follows:
  + **MBR**: Master Boot Record (MBR).
  + **GPT**: GUID Partition Table (GPT).

  -> This parameter is mandatory for source server registration.

* `os_disk` - (Optional, Bool) Specifies whether the disk is the system disk.

* `relation_name` - (Optional, String) Specifies the name of the paired target server disk in Linux.

* `inode_size` - (Optional, List) Specifies the number of inodes.

* `id` - (Optional, Int) Specifies the disk ID.

* `adjust_size` - (Optional, Int) Specifies the new size.

* `need_migration` - (Optional, Bool) Specifies whether the volume needs to be migrated.

<a name="disks_physical_volumes_struct"></a>
The `physical_volumes` block supports:

* `device_use` - (Optional, String) Specifies the partition function.
  Values can be general, boot, or OS partition.

* `file_system` - (Optional, String) Specifies the file system type.

* `index` - (Optional, Int) Specifies the serial number.

* `mount_point` - (Optional, String) Specifies the mount point.

* `name` - (Optional, String) Specifies the volume name. In Windows, it indicates the drive letter, and in Linux, it
  indicates the device ID.

* `size` - (Optional, Int) Specifies the size.

* `used_size` - (Optional, Int) Specifies the used space.

* `inode_size` - (Optional, Int) Specifies the number of inodes.

* `inode_nums` - (Optional, Int) Specifies the number of inode nodes.

* `uuid` - (Optional, String) Specifies the GUID.

* `size_per_cluster` - (Optional, Int) Specifies the size of each cluster.

* `id` - (Optional, Int) Specifies the database record ID.

* `adjust_size` - (Optional, Int) Specifies the new size.

* `need_migration` - (Optional, Bool) Specifies whether the volume needs to be migrated.

<a name="btrfs_list_struct"></a>
The `btrfs_list` block supports:

* `name` - (Required, String) Specifies the file system name.

* `label` - (Required, String) Specifies the file system tag. If no tag exists, the value is an empty string.

* `uuid` - (Required, String) Specifies the UUID of the file system.

* `device` - (Required, String) Specifies the device names of the Btrfs file system.

* `size` - (Required, Int) Specifies the space occupied by the file system.

* `nodesize` - (Required, Int) Specifies the Btrfs node size.

* `sectorsize` - (Required, Int) Specifies the sector size.

* `data_profile` - (Required, String) Specifies the data profile (RAD).

* `system_profile` - (Required, String) Specifies the file system profile (RAD).

* `metadata_profile` - (Required, String) Specifies the metadata profile (RAD).

* `global_reserve1` - (Required, String) Specifies the Btrfs file system information.

* `g_vol_used_size` - (Required, Int) Specifies the used space of the Btrfs volume.

* `default_subvolid` - (Required, String) Specifies the ID of the default subvolumn.

* `default_subvol_name` - (Required, String) Specifies the name of the default subvolumn.

* `default_subvol_mountpath` - (Required, String) Specifies the mount path of the default subvolumn or Btrfs file system.

* `subvolumn` - (Required, List) Specifies the subvolumn information.
  The [subvolumn](#btrfs_list_subvolumn_struct) structure is documented below.

<a name="btrfs_list_subvolumn_struct"></a>
The `subvolumn` block supports:

* `uuid` - (Required, String) Specifies the UUID of the parent volume.

* `is_snapshot` - (Required, String) Specifies whether the subvolumn is a snapshot.

* `subvol_id` - (Required, String) Specifies the subvolumn ID.

* `parent_id` - (Required, String) Specifies the parent volume ID.

* `subvol_name` - (Required, String) Specifies the subvolumn name.

* `subvol_mount_path` - (Required, String) Specifies the mount path of the subvolumn.

<a name="networks_struct"></a>
The `networks` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the NIC name.

* `ip` - (Required, String, NonUpdatable) Specifies the IP address bound to the NIC.

* `netmask` - (Required, String, NonUpdatable) Specifies the subnet mask.

* `gateway` - (Required, String, NonUpdatable) Specifies the gateway.

* `mac` - (Required, String, NonUpdatable) Specifies the MAC address.

* `ipv6` - (Optional, String, NonUpdatable) Specifies the IPv6 address.

* `mtu` - (Optional, Int, NonUpdatable) Specifies the NIC MTU.

  -> This parameter is mandatory for Linux.

* `id` - (Optional, String) Specifies the database record ID.

<a name="volume_groups_struct"></a>
The `volume_groups` block supports:

* `components` - (Optional, String) Specifies the physical volume information.

* `free_size` - (Optional, Int) Specifies the available space.

* `logical_volumes` - (Optional, List) Specifies the logical volume information.
  The [logical_volumes](#volume_groups_logical_volumes_struct) structure is documented below.

* `name` - (Optional, String) Specifies the name.

* `size` - (Optional, Int) Specifies the size.

* `id` - (Optional, Int) Specifies the volume group ID.

* `adjust_size` - (Optional, Int) Specifies the new size.

* `need_migration` - (Optional, Bool) Specifies whether the volume needs to be migrated.

<a name="volume_groups_logical_volumes_struct"></a>
The `logical_volumes` block supports:

* `file_system` - (Required, String) Specifies the file system.

* `inode_size` - (Required, Int) Specifies the number of inodes.

* `mount_point` - (Required, String) Specifies the mount point.

* `name` - (Required, String) Specifies the name.

* `size` - (Required, String) Specifies the size.

* `used_size` - (Required, String) Specifies the used space.

* `free_size` - (Required, Int) Specifies the available space.

* `block_count` - (Optional, Int) Specifies the number of blocks.

* `block_size` - (Optional, Int) Specifies the block size.

* `inode_nums` - (Optional, String) Specifies the number of inodes.

* `device_use` - (Optional, String) Specifies the partition function.
  Values can be general, boot, or OS partition.

* `id` - (Optional, Int) Specifies the logical volume ID.

* `adjust_size` - (Optional, Int) Specifies the new size.

* `need_migration` - (Optional, Bool) Specifies whether the volume needs to be migrated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `add_date` - Indicates the time when the source server was registered.

* `connected` - Indicates whether the Agent installed on the source server is connected to SMS.

* `init_target_server` - Indicates the recommended configuration for the target server.
  The [init_target_server](#init_target_server_struct) structure is documented below.

* `current_task` - Indicates the migration task associated with the source server.
  The [current_task](#current_task_struct) structure is documented below.

* `checks` - Indicates the environment check information for the source server.
  The [checks](#checks_struct) structure is documented below.

* `state_action_time` - Indicates the timestamp when the status of the source server last changed. The source server
  status is defined by `state`.

* `replicatesize` - Indicates the volume of data that has been migrated, the unit is bytes.

* `totalsize` - Indicates the volume of data to be migrated, the unit is bytes.

* `last_visit_time` - Indicates the timestamp when the Agent connection status last changed.

* `stage_action_time` - Indicates the timestamp when the migration stage of the source server last changed. The migration
  stage is defined by `migration_cycle`.

* `adjust_disk` - Indicates whether the disk is adjusted.

<a name="init_target_server_struct"></a>
The `init_target_server` block supports:

* `disks` - Indicates the information about the recommended target server disks.
  The [disks](#init_target_server_disks_struct) structure is documented below.

* `volume_groups` - Indicates the volume groups.
  The [volume_groups](#init_target_server_volume_groups_struct) structure is documented below.

  -> This parameter is mandatory for Linux.

<a name="init_target_server_disks_struct"></a>
The `disks` block supports:

* `name` - Indicates the disk name.

* `size` - Indicates the disk size, the unit is bytes.

* `device_use` - Indicates the disk function.
  Values can be as follows:
  + **BOOT**: Boot device.
  + **OS**: System device.
  + **NORMAL**: General device.

* `used_size` - Indicates the used disk space, the unit is bytes.

* `id` - Indicates the disk ID.

* `adjust_size` - Indicates the new size.

* `need_migration` - Indicates whether the volume needs to be migrated.

* `physical_volumes` - Indicates the enumeration list.
  The [physical_volumes](#init_target_server_disks_physical_volumes_struct) structure is documented below.

<a name="init_target_server_disks_physical_volumes_struct"></a>
The `physical_volumes` block supports:

* `device_use` - Indicates the partition function. The partition can be a general, boot, or OS partition.

* `file_system` - Indicates the file system type.

* `index` - Indicates the serial number.

* `mount_point` - Indicates the mount point.

* `name` - Indicates the volume name. In Windows, it indicates the drive letter, and in Linux, it indicates the device ID.

* `size` - Indicates the size.

* `inode_size` - Indicates the number of inodes.

* `used_size` - Indicates the used space.

* `uuid` - Indicates the GUID, which can be obtained from the source server.

* `id` - Indicates the database record ID.

* `adjust_size` - Indicates the new size.

* `need_migration` - Indicates whether the volume needs to be migrated.

<a name="init_target_server_volume_groups_struct"></a>
The `volume_groups` block supports:

* `components` - Indicates the physical volume information.

* `free_size` - Indicates the available space.

* `logical_volumes` - Indicates the logical volume information.
  The [logical_volumes](#init_target_server_volume_groups_logical_volumes_struct) structure is documented below.

* `name` - Indicates the name.

* `size` - Indicates the size.

* `id` - Indicates the volume group ID.

* `adjust_size` - Indicates the new size.

* `need_migration` - Indicates whether the volume needs to be migrated.

<a name="init_target_server_volume_groups_logical_volumes_struct"></a>
The `logical_volumes` block supports:

* `block_count` - Indicates the number of blocks.

* `block_size` - Indicates the block size.

* `file_system` - Indicates the file system.

* `inode_size` - Indicates the number of inodes.

* `inode_nums` - Indicates the number of inode nodes.

* `device_use` - Indicates the partition function. The partition can be a general, boot, or OS partition.

* `mount_point` - Indicates the mount point.

* `name` - Indicates the name.

* `size` - Indicates the size.

* `used_size` - Indicates the used space.

* `free_size` - Indicates the available space.

* `id` - Indicates the logical volume ID.

* `adjust_size` - Indicates the new size.

* `need_migration` - Indicates whether the volume needs to be migrated.

<a name="current_task_struct"></a>
The `current_task` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `type` - Indicates the task type.

* `state` - Indicates the task status.

* `start_date` - Indicates the start time.

* `speed_limit` - Indicates the migration rate limit.

* `migrate_speed` - Indicates the migration rate.

* `start_target_server` - Indicates whether the target server is started.

* `vm_template_id` - Indicates the server template ID.

* `region_id` - Indicates the region ID.

* `project_name` - Indicates the project name.

* `project_id` - Indicates the project ID.

* `target_server` - Indicates the information about the target server.
  The [target_server](#current_task_target_server_struct) structure is documented below.

* `log_collect_status` - Indicates the log collection status.

* `exist_server` - Indicates whether an existing server is used as the target server.

* `use_public_ip` - Indicates whether a public IP address is used for migration.

* `clone_server` - Indicates the information about the cloned server.
  The [clone_server](#current_task_clone_server_struct) structure is documented below.

<a name="current_task_target_server_struct"></a>
The `target_server` block supports:

* `vm_id` - Indicates the ID of the target server.

* `name` - Indicates the name of the target server.

<a name="current_task_clone_server_struct"></a>
The `clone_server` block supports:

* `vm_id` - Indicates the ID of the cloned server.

* `name` - Indicates the name of the cloned server.

* `clone_error` - Indicates the error returned for a clone failure.

* `clone_state` - Indicates the clone status.

* `error_msg` - Indicates the error message returned for a clone failure.

<a name="checks_struct"></a>
The `checks` block supports:

* `id` - Indicates the check item ID.

* `params` - Indicates the parameters.

* `name` - Indicates the check item name.

* `result` - Indicates the check result.
  Values can be as follows:
  + **OK**: The check is passed.
  + **WARN**: A warning is generated.
  + **ERROR**: The check fails.

* `error_code` - Indicates the returned error code.

* `error_or_warn` - Indicates the returned error or warning.

* `error_params` - Indicates the parameters that failed the check.

## Import

SMS source servers can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_sms_source_server.demo <id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `virtualization_type`, `linux_block_check`, `domain_id`, `has_rsync`,
`paravirtualization`, `raw_devices`, `driver_files`, `system_services`, `account_rights`, `boot_loader`, `system_dir`，
`kernel_version`, `start_type`, `io_read_wait`, `platform`, `migprojectid`, `copystate`.
It is generally recommended running `terraform plan` after importing a source server.
You can then decide if changes should be applied to the scan task, or the resource definition should be updated to align
with the source server. Also you can ignore changes as below.

```hcl
resource "huaweicloud_sms_source_server" "test" {
  ...

  lifecycle {
    ignore_changes = [
      virtualization_type, linux_block_check, domain_id, has_rsync, paravirtualization, raw_devices, driver_files,
      system_services, account_rights, boot_loader, system_dir, kernel_version, start_type, io_read_wait, platform,
      migprojectid, copystate,
    ]
  }
}
```
