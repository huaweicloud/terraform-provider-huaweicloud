---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo"
description: |-
  Manages a SFS Turbo resource within HuaweiCloud.
---

# huaweicloud_sfs_turbo

Manages a SFS Turbo resource within HuaweiCloud.

## Example Usage

### Create a STANDARD Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "sfs-turbo-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create an HPC Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "sfs-turbo-1"
  size              = 3686
  share_proto       = "NFS"
  share_type        = "HPC"
  hpc_bandwidth     = "40M"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az
}
```

### Create an HPC CACHE Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "huaweicloud_sfs_turbo" "test" {
  name                = "sfs-turbo-1"
  size                = 4096
  share_proto         = "NFS"
  share_type          = "HPC_CACHE"
  hpc_cache_bandwidth = "2G"
  vpc_id              = var.vpc_id
  subnet_id           = var.subnet_id
  security_group_id   = var.secgroup_id
  availability_zone   = var.test_az
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo resource.

* `name` - (Required, String) Specifies the name of an SFS Turbo file system. The value contains `4` to `64`
  characters and must start with a letter.

* `size` - (Required, Int) Specifies the capacity of a sharing file system, in GB.
  + If `share_type` is set to **STANDARD** or **PERFORMANCE**, the value ranges from `500` to `32,768`, and ranges from
  `10,240` to `327,680` for an enhanced file system.

  + If `share_type` is set to **HPC**, the value ranges from `3,686` to `1,048,576` when `hpc_bandwidth` is set to
  **20M**, and ranges from `1,228` to `1,048,576` when `hpc_bandwidth` is set to **40M**, **125M**, **250M**, **500M**
  or **1000M**. The capacity must be a multiple of 1.2TiB, which needs to be rounded down after converting to GiB.
  Such as 3.6TiB->3686GiB, 4.8TiB->4915GiB, 8.4TiB->8601GiB.

  + If `share_type` is set to **HPC_CACHE**, the value ranges from `4,096` to `1,048,576`, and the step size is `1,024`.
  The minimum capacity(GB) should be equal to `2,048` multiplying the HPC cache bandwidth size(GB/s).
  Such as the minimum capacity is `4,096` when `hpc_cache_bandwidth` is set to **2G**, the minimum capacity is `8,192`
  when `hpc_cache_bandwidth` is set to **4G**, the minimum capacity is `16,384` when `hpc_cache_bandwidth` is set to
  **8G**.

  -> The file system capacity can only be expanded, not reduced.

* `share_proto` - (Optional, String, ForceNew) Specifies the protocol for sharing file systems. The valid value is NFS.
  Changing this will create a new resource.

* `share_type` - (Optional, String, ForceNew) Specifies the file system type. Changing this will create a new resource.
  Valid values are **STANDARD**, **PERFORMANCE**, **HPC** and **HPC_CACHE**.
  Defaults to **STANDARD**.

  -> The share type **HPC_CACHE** only support in postpaid charging mode.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone where the file system is located.
  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of the subnet. Changing this will create a new
  resource.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `enhanced` - (Optional, Bool, ForceNew) Specifies whether the file system is enhanced or not. Changing this will
  create a new resource.

  This parameter is valid only when `share_type` is set to **STANDARD** or **PERFORMANCE**.

* `hpc_bandwidth` - (Optional, String, ForceNew) Specifies the HPC bandwidth. Changing this will create a new resource.
  This parameter is valid and required when `share_type` is set to **HPC**.
  Valid values are: **20M**, **40M**, **125M**, **250M**, **500M** and **1000M**.

* `hpc_cache_bandwidth` - (Optional, String) Specifies the HPC cache bandwidth(GB/s).
  This parameter is valid and required when `share_type` is set to **HPC_CACHE**.
  Valid values are: **2G**, **4G**, **8G**, **16G**, **24G**, **32G** and **48G**.

* `crypt_key_id` - (Optional, String, ForceNew) Specifies the ID of a KMS key to encrypt the file system. Changing this
  will create a new resource.

* `dedicated_flavor` - (Optional, String, ForceNew) Specifies the VM flavor used for creating a dedicated file system.

* `dedicated_storage_id` - (Optional, String, ForceNew) Specifies the ID of the dedicated distributed storage used
  when creating a dedicated file system.

* `auto_create_security_group_rules` - (Optional, String) Specifies whether to automatically create security
  group rules. **true** means automatically create security group rules.
  **false** means not automatically create security group rules. Defaults to **true**.
  This field cannot be edited individually. Editing this field alone will not make any changes to the resource.
  Editing this field will only take effect when the `security_group_id` field is changed.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the file system. Changing this
  will create a new resource.

* `backup_id` - (Optional, String, ForceNew) Specifies the backup ID.

  -> This parameter is mandatory when a file system is created from a backup.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the SFS Turbo.

-> **NOTE:**
SFS Turbo will create two private IP addresses and one virtual IP address under the subnet you specified. To ensure
normal use, SFS Turbo will enable the inbound rules for ports *111*, *445*, *2049*, *2051*, *2052*, and *20048* in the
security group you specified.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the SFS Turbo.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.
  Changing this parameter will create a new cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the SFS Turbo.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the SFS Turbo.
  If `period_unit` is set to **month**, the value ranges from `1` to `11`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.  
  The valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `status` - The status of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - The mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 10 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sfs_turbo.test <id>
```

Note that the imported state may not be identical to your resource definition, due to payment attributes missing from
the API response.
The missing attributes include: `charging_mode`, `period_unit`, `period`, `auto_renew`,
`auto_create_security_group_rules`, `dedicated_flavor`, `dedicated_storage_id`.
It is generally recommended running `terraform plan` after importing an instance.
You can ignore changes as below.

```hcl
resource "huaweicloud_sfs_turbo" "test" {
  ...

  lifecycle {
    ignore_changes = [
      charging_mode, period_unit, period, auto_renew,
      auto_create_security_group_rules, dedicated_flavor, dedicated_storage_id,
    ]
  }
}
```
