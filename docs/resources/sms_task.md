---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task"
description: ""
---

# huaweicloud_sms_task

Manages an SMS migration task resource within HuaweiCloud.

## Example Usage

```hcl
variable "source_server" {}
variable "template_id" {}

resource "huaweicloud_sms_task" "migration" {
  type             = "MIGRATE_FILE"
  os_type          = "LINUX"
  source_server_id = var.source_server
  vm_template_id   = var.template_id
  action           = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the target server is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the migration task. Available values are
  **MIGRATE_FILE**(file-level migration) and **MIGRATE_BLOCK**(block-level migration).
  Changing this parameter will create a new resource.

  + For Linux servers, SMS supports block-level and file-level migrations. Block-level migration has
    high efficiency but poor compatibility while file-level migration has low efficiency but excellent compatibility.
  + For Windows servers, SMS only supports highly efficient block-level migration.

* `os_type` - (Required, String, ForceNew) Specifies the OS type of the source server. The value can be **WINDOWS** and **LINUX**.
  Changing this parameter will create a new resource.

* `source_server_id` - (Required, String, ForceNew) Specifies the ID of the source server.
  Changing this parameter will create a new resource.

* `vm_template_id` - (Optional, String, ForceNew) Specifies the template used to create the target server automatically.
   This parameter and `target_server_id` are alternative. Changing this parameter will create a new resource.

* `target_server_id` - (Optional, String, ForceNew) Specifies the existing server ID as the target server.
   This parameter and `vm_template_id` are alternative. Changing this parameter will create a new resource.

* `target_server_disks` - (Optional, List, ForceNew) Specifies the disk configurations of the target server.
  If omitted, it will be obtained from the source server. The [object](#target_server_disks_object)
  is documented below. Changing this parameter will create a new resource.

* `use_public_ip` - (Optional, Bool, ForceNew) Specifies whether to use a public IP address for migration.
  The default value is `true`. Changing this parameter will create a new resource.

* `migration_ip` - (Optional, String, ForceNew) Specifies the IP address of the target server.
  Use the EIP of the target server if the migration network type is Internet.
  Use the private IP address of the target server if the migration network type is Direct Connect or VPN.
  Changing this parameter will create a new resource.

* `start_target_server` - (Optional, Bool, ForceNew) Specifies whether to start the target server after the migration.
  The default value is `true`. Changing this parameter will create a new resource.

* `syncing` - (Optional, Bool, ForceNew) - Specifies whether to perform a continuous synchronization after the first replication.
  The default value is `false`. Changing this parameter will create a new resource.

* `action` - (Optional, String) Specifies the operation after the task is created.
  The value can be **start**, **stop** and **restart**.

* `speed_limit` - (Optional, List) Specifies the information about the time period-based rate limiting rules.
  The [speed_limit](#speed_limit_struct) structure is documented below.

* `configurations` - (Optional, List) Specifies the settings of advanced migration options.
  The [configurations](#configurations_struct) structure is documented below.

* `project_id` - (Optional, String, ForceNew) Specifies the project ID where the target server is located.
  If omitted, the default project in the region will be used. Changing this parameter will create a new resource.

<a name="speed_limit_struct"></a>
The `speed_limit` block supports:

* `start` - (Required, String) Specifies the start time of a period. The format is **XX:XX**.

* `end` - (Required, String) Specifies the end time of a period. The format is **XX:XX**.

* `speed` - (Required, Int) Specifies the migration rate limit for the specified period, the unit is Mbit/s.

* `over_speed_threshold` - (Optional, Float) Specifies the overspeed threshold for stopping migration, the unit is percentage.

<a name="configurations_struct"></a>
The `configurations` block supports:

* `config_key` - (Required, String) Specifies the advanced migration option.
  The value can be **EXCLUDE_MIGRATE_PATH**, **SYNC_EXCLUDE_PATH**, **ONLY_SYNC_PATH** and so on.

* `config_value` - (Required, String) Specifies the value specified for the advanced migration option.

* `config_status` - (Optional, String) Specifies the reserved field that describes the configuration status.

<a name="target_server_disks_object"></a>
The `target_server_disks` block supports:

* `name` - (Required, String, ForceNew) Specifies the disk name, e.g. "/dev/sda".
  Changing this parameter will create a new resource.

* `size` - (Required, Int, ForceNew) Specifies the volume size in MB. Changing this parameter will create a new resource.

* `device_type` - (Required, String, ForceNew) Specifies the disk type. The value can be **NORMAL** and **BOOT**.
  Changing this parameter will create a new resource.

* `disk_id` - (Required, String, ForceNew) Specifies the disk index, e.g. "0".
  Changing this parameter will create a new resource.

* `used_size` - (Optional, Int, ForceNew) Specifies the used space in MB. Changing this parameter will create a new resource.

* `physical_volumes` - (Optional, List, ForceNew) Specifies an array of physical volume information.
  The [object](#physical_volumes_object) is documented below. Changing this parameter will create a new resource.

<a name="physical_volumes_object"></a>
The `physical_volumes` block supports:

* `name` - (Required, String, ForceNew) Specifies the volume name. In Windows, it indicates the drive letter,
  and in Linux, it indicates the device ID, e.g. "/dev/sda1".
  Changing this parameter will create a new resource.

* `size` - (Required, Int, ForceNew) Specifies the volume size in MB. Changing this parameter will create a new resource.

* `device_type` - (Required, String, ForceNew) Specifies the partition type. The value can be **NORMAL** and **OS**.
  Changing this parameter will create a new resource.

* `file_system` - (Required, String, ForceNew) Specifies the file system type, e.g. "ext4".
  Changing this parameter will create a new resource.

* `mount_point` - (Required, String, ForceNew) Specifies the mount point, e.g. "/".
  Changing this parameter will create a new resource.

* `index` - (Required, Int, ForceNew) Specifies the serial number of the volume.
  Changing this parameter will create a new resource.

* `used_size` - (Optional, Int, ForceNew) Specifies the used space in MB.
  Changing this parameter will create a new resource.

* `uuid` - (Optional, String, ForceNew) Specifies the GUID of the volume.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `state` - Indicates the status of the migration task.

* `enterprise_project_id` - Indicates the enterprise project id of the target server.

* `target_server_name` - Indicates the name of the target server.

* `migrate_speed` - Indicates the migration rate, in MB/s.

* `passphrase` - Indicates the certificate passphrase of the secure transmission channel.

* `migrate_type` - Indicates the migration method.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

## Import

SMS migration tasks can be imported by `id`, e.g.

```sh
terraform import huaweicloud_sms_task.demo 6402c49b-7d9a-413e-8b5f-a7307f7d5679
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `use_public_ip`, `syncing` and `action`.
It is generally recommended running `terraform plan` after importing a migration task.
You can then decide if changes should be applied to the task, or the resource definition should be
updated to align with the task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_sms_task" "demo" {
    ...

  lifecycle {
    ignore_changes = [
      use_public_ip, syncing, action,
    ]
  }
}
```
