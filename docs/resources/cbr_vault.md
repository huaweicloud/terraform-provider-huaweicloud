---
subcategory: "Cloud Backup and Recovery (CBR)"
---

# huaweicloud_cbr_vault

Manages a CBR Vault resource within Huaweicloud.

## Example Usage

### Create a server type vault

```hcl
variable "vault_name" {}
variable "ecs_instance_id" {}
variable "evs_volume_id" {}

data "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "server"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 100

  resources {
    server_id = var.ecs_instance_id
  
    excludes = [
      var.evs_volume_id
    ]
  }

  tags = {
    foo = "bar"
  }
}
```

### Create a disk type vault

```hcl
variable "vault_name" {}
variable "evs_volume_id" {}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  type             = "disk"
  protection_type  = "backup"
  consistent_level = "crash_consistent"
  size             = 50
  auto_expand      = true

  resources {
    includes = [
      var.evs_volume_id
    ]
  }

  tags = {
    foo = "bar"
  }
}
```

### Create an SFS turbo type vault

```hcl
variable "vault_name" {}
variable "sfs_turbo_id" {}

resource "huaweicloud_cbr_vault" "test" {
  name             = var.vault_name
  consistent_level = "crash_consistent"
  type             = "turbo"
  protection_type  = "backup"
  size             = 1000

  resources {
    includes = [
      var.sfs_turbo_id
    ]
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
  name             = var.vault_name
  consistent_level = "crash_consistent"
  type             = "turbo"
  protection_type  = "replication"
  size             = 1000
}
```

## Argument reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CBR vault. If omitted, the
  provider-level region will be used. Changing this will create a new vault.

* `name` - (Required, String) Specifies a unique name of the CBR vault. This parameter can contain a maximum of 64
  characters, which may consist of letters, digits, underscores(_) and hyphens (-).

* `type` - (Required, String, ForceNew) Specifies the object type of the CBR vault.
  Changing this will create a new vault. Vaild values are as follows:
  + **server** (Cloud Servers)
  + **disk** (EVS Disks)
  + **turbo** (SFS Turbo file systems)

* `consistent_level` - (Required, String, ForceNew) Specifies the backup specifications.
  The valid values are as follows:
  + **[crash_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**
  + **[app_consistent](https://support.huaweicloud.com/intl/en-us/usermanual-cbr/cbr_03_0109.html)**

  Only server type vaults support application consistent. Changing this will create a new vault.

* `protection_type` - (Required, String, ForceNew) Specifies the protection type of the CBR vault.
  The valid values are **backup** and **replication**. Vaults of type **disk** don't support **replication**.
  Changing this will create a new vault.

* `size` - (Required, Int) Specifies the vault sapacity, in GB. The valid value range is `1` to `10,485,760`.

* `auto_expand` - (Optional, Bool) Specifies to enable auto capacity expansion for the backup protection type vault.
  Default to **false**.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies a unique ID in UUID format of enterprise project.
  Changing this will create a new vault.

* `policy_id` - (Optional, String) Specifies a policy to associate with the CBR vault.
  `policy_id` cannot be used with the vault of replicate protection type.

* `resources` - (Optional, List) Specifies an array of one or more resources to attach to the CBR vault.
  The [object](#cbr_vault_resources) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CBR vault.

<a name="cbr_vault_resources"></a>
The `resources` block supports:

* `server_id` - (Optional, String) Specifies the ID of the ECS instance to be backed up.

* `excludes` - (Optional, List) Specifies the array of disk IDs which will be excluded in the backup.
  Only **server** vault support this parameter.

* `includes` - (Optional, List) Specifies the array of disk or SFS file system IDs which will be included in the backup.
  Only **disk** and **turbo** vault support this parameter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `allocated` - The allocated capacity of the vault, in GB.

* `used` - The used capacity, in GB.

* `spec_code` - The specification code.

* `status` - The vault status.

* `storage` - The name of the bucket for the vault.

## Import

Vaults can be imported by their `id`. For example,

```
$ terraform import huaweicloud_cbr_vault.test 01c33779-7c83-4182-8b6b-24a671fcedf8
```
