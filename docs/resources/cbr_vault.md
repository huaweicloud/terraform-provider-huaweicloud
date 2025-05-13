---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_vault"
description: ""
---

# huaweicloud_cbr_vault

Manages a CBR vault resource within HuaweiCloud.

## Example Usage

### Create a server type vault

```hcl
variable "vault_name" {}
variable "ecs_instance_id" {}
variable "attached_volume_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "server"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 100

  resources {
    server_id = var.ecs_instance_id
    excludes  = var.attached_volume_ids
  }

  tags = {
    foo = "bar"
  }
}
```

### Create a server type vault and associate backup and reprecation policies

```hcl
variable "destination_region" {}
variable "destination_vault_name" {}
variable "vault_name" {}
variable "backup_policy_id" {}
variable "replication_policy_id" {}

resource "huaweicloud_cbr_vault" "destination" {
  region          = var.destination_region
  name            = var.destination_vault_name
  type            = "server"
  protection_type = "replication"
  size            = 500
}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "server"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 500

  auto_bind  = true
  bind_rules = {
    # {auto_bind_key} = {auto_bind_value}
    service_name = "xxx"
  }

  policy {
    id = var.backup_policy_id
  }
  policy {
    id                   = var.replication_policy_id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
```

### Create a disk type vault

```hcl
variable "vault_name" {}
variable "evs_volume_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_vault" "test" {
  name            = var.vault_name
  type            = "disk"
  protection_type = "backup"
  size            = 50
  auto_expand     = true
  locked          = false

  resources {
    includes = var.evs_volume_ids
  }

  tags = {
    foo = "bar"
  }
}
```

### Create an SFS turbo type vault

```hcl
variable "vault_name" {}
variable "sfs_turbo_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_vault" "test" {
  name            = var.vault_name
  type            = "turbo"
  protection_type = "backup"
  size            = 1000

  resources {
    includes = var.sfs_turbo_ids
  }

  tags = {
    foo = "bar"
  }
}
```

### Create an SFS turbo type vault with replicate protection type

```hcl
variable "vault_name" {}

resource "huaweicloud_cbr_vault" "test" {
  name            = var.vault_name
  type            = "turbo"
  protection_type = "replication"
  size            = 1000
}
```

### Create a Workspace type vault

```hcl
variable "vault_name" {}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "workspace"
  protection_type  = "backup"
  size             = 100
  consistent_level = "crash_consistent"
}
```

### Create a VMware type vault

```hcl
variable "vault_name" {}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "vmware"
  protection_type  = "backup"
  size             = 100
  consistent_level = "crash_consistent"
}
```

### Create a file type vault

```hcl
variable "vault_name" {}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "file"
  protection_type  = "backup"
  size             = 100
  consistent_level = "crash_consistent"
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CBR vault. If omitted, the
  provider-level region will be used. Changing this will create a new vault.

* `cloud_type` - (Required, String, ForceNew) Specifies the cloud type of the vault.  
  Changing this will create a new vault.

* `name` - (Required, String) Specifies a unique name of the CBR vault. This parameter can contain a maximum of 64
  characters, which may consist of letters, digits, underscores(_) and hyphens (-).

* `type` - (Required, String, ForceNew) Specifies the object type of the CBR vault.
  Changing this will create a new vault. Valid values are as follows:
  + **server** (Elastic Cloud Server)
  + **disk** (EVS Disk)
  + **turbo** (SFS Turbo file system)
  + **workspace** (Workspace Desktop)
  + **vmware** (VMware)
  + **file** (File System)

* `protection_type` - (Required, String, ForceNew) Specifies the protection type of the CBR vault.
  The valid values are **backup** and **replication**. Vaults of type **disk** don't support **replication**.
  Changing this will create a new vault.

* `size` - (Required, Int) Specifies the vault capacity, in GB. The valid value range is `1` to `10,485,760`.

  -> You cannot update `size` if the vault is **prePaid** mode.

* `consistent_level` - (Optional, String) Specifies the consistent level (specification) of the vault.
  The valid values are as follows:
  + **[crash_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**
  + **[app_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**

  Only **server** type vaults support application consistent and defaults to **crash_consistent**, and only
  **crash_consistent** can be updated to **app_consistent**.  
  The **workspace** type vaults does not support application consistent.

* `auto_expand` - (Optional, Bool) Specifies to enable auto capacity expansion for the backup protection type vault.
  Defaults to **false**.

  -> You cannot configure `auto_expand` if the vault is **prePaid** mode.

* `auto_expand` - (Optional, Bool) Specifies Locked status of the vault. Defaults to **false**.

* `auto_bind` - (Optional, Bool) Specifies whether automatic association is enabled. Defaults to **false**.

* `bind_rules` - (Optional, Map) Specifies the tags to filter resources for automatic association with **auto_bind**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the vault belongs.

* `policy` - (Optional, List) Specifies the policy details to associate with the CBR vault.
  The [object](#cbr_vault_policies) structure is documented below.

* `resources` - (Optional, List) Specifies an array of one or more resources to attach to the CBR vault.  
  This feature is not supported for the **vmware** type and the **file** type.  
  The [object](#cbr_vault_resources) structure is documented below.

* `backup_name_prefix` - (Optional, String, ForceNew) Specifies the backup name prefix.
  Changing this will create a new vault.

-> If configured, the names of all automatic backups generated for the vault will use this prefix.

* `is_multi_az` - (Optional, Bool, ForceNew) Specifies whether multiple availability zones are used for backing up.
  Defaults to **false**.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CBR vault.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the vault.
  The valid values are as follows:
  + **prePaid**: the yearly/monthly billing mode.
  + **postPaid**: the pay-per-use billing mode.

  Changing this will create a new vault.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the vault.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this will create a new vault.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the vault.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 5.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this will create a new vault.

* `auto_renew` - (Optional, String) Specifies whether auto renew is enabled.
  Valid values are **true** and **false**. Defaults to **false**.

<a name="cbr_vault_policies"></a>
The `policy` block supports:

* `id` - (Required, String) Specifies the policy ID.

* `destination_vault_id` - (Optional, String) Specifies the ID of destination vault to which the replication policy
  will associated.

-> Only one policy of each type (backup and replication) can be associated.

<a name="cbr_vault_resources"></a>
The `resources` block supports:

* `server_id` - (Optional, String) Specifies the ID of the ECS instance to be backed up.

* `excludes` - (Optional, List) Specifies the array of disk IDs which will be excluded in the backup.
  Only **server** vault support this parameter.

* `includes` - (Optional, List) Specifies the array of disk or SFS file system IDs which will be included in the backup.
  Only **disk** and **turbo** vault support this parameter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `allocated` - The allocated capacity of the vault, in GB.

* `used` - The used capacity, in GB.

* `spec_code` - The specification code.

* `status` - The vault status.

* `storage` - The name of the bucket for the vault.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

Vaults can be imported by their `id`. For example,

```bash
$ terraform import huaweicloud_cbr_vault.test 01c33779-7c83-4182-8b6b-24a671fcedf8
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `period_unit`, `period`, `auto_renew`.
It is generally recommended running `terraform plan` after importing a vault.
You can then decide if changes should be applied to the vault, or the resource definition should be updated to align
with the vault. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cbr_vault" "test" {
  ...

  lifecycle {
    ignore_changes = [
      period_unit, period, auto_renew,
    ]
  }
}
```
