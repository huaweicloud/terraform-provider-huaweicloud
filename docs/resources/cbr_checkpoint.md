---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_checkpoint"
description: ""
---

# huaweicloud_cbr_checkpoint

Using this resource to manage a checkpoint and related resource backups within HuaweiCloud.

## Example Usage

```hcl
variable "vault_id" {}
variable "checkpoint_name" {}
variable "backup_resource_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_checkpoint" "test" {
  vault_id = var.vault_id
  name     = var.checkpoint_name

  dynamic "backups" {
    for_each = var.backup_resource_ids

    content {
      type        = "OS::Nova::Server"
      resource_id = backups.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the vault and backup resources are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `vault_id` - (Required, String, ForceNew) Specifies the ID of the vault where the checkpoint to create.  
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the checkpoint.  
  The valid length is limited from `1` to `64`, only Chinese and English letters, digits, hyphens (-) and
  underscores (_) are allowed.  
  Changing this will create a new resource.

* `backups` - (Required, List, ForceNew) Specifies the backup resources configuration.
  The [backups](#cbr_checkpoint_backups_args) structure is documented below.  
  Changing this will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the checkpoint.  
  The valid length is limited from `0` to `255`.
  Changing this will create a new resource.

* `incremental` - (Optional, Bool, ForceNew) Specifies whether the backups are incremental backups.  
  Defaults to **false**. Changing this will create a new resource.

<a name="cbr_checkpoint_backups_args"></a>
The `backups` block supports:

* `type` - (Required, String, ForceNew) Specifies the type of the backup resource.  
  The valid values are as follows:
  + **OS::Nova::Server**
  + **OS::Cinder::Volume**
  + **OS::Ironic::BareMetalServer**
  + **OS::Native::Server**
  + **OS::Sfs::Turbo**
  + **OS::Workspace::DesktopV2**

  Changing this will create a new resource.

* `resource_id` - (Required, String, ForceNew) Specifies the ID of backup resource.  
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `backups` - The backup resources configuration.
  The [backups](#cbr_checkpoint_backup_attr) structure is documented below.  

* `created_at` - The creation time of the checkpoint.

* `status` - The status of the checkpoint.

<a name="cbr_checkpoint_backup_attr"></a>
The `backups` block supports:

* `id` - The backup ID.

* `resource_size` - The backup resource size.

* `status` - The backup status.

* `protected_at` - The backup time.

* `updated_at` - The latest update time of the backup.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 20 minutes.
