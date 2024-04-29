---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backup_share"
description: ""
---

# huaweicloud_cbr_backup_share

Using this resource to share backups with other members within HuaweiCloud (in the same region).

-> Currently, only Server backup type support to manage shared members.
   And a backup can only create one of this resource.

## Example Usage

```hcl
variable "backup_id" {}
variable "shared_project_ids" {
  type = list(string)
}

resource "huaweicloud_cbr_backup_share" "test" {
  backup_id = var.backup_id

  dynamic "members" {
    for_each = var.shared_project_ids

    content {
      dest_project_id = members.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the backup and sharing project belong.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `backup_id` - (Required, String, ForceNew) Specifies the backup ID.  
  Changing this will create a new resource.

* `members` - (Required, List) Specifies the list of shared members configuration.
  The [members](#cbr_backup_share_members_args) structure is documented below.  

<a name="cbr_backup_share_members_args"></a>
The `members` block supports:

* `dest_project_id` - (Required, String) Specifies the ID of the project with which the backup is shared.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also backup ID) in UUID format.

* `members` - The list of shared members configuration.
  The [members](#cbr_checkpoint_backup_attr) structure is documented below.  

<a name="cbr_checkpoint_backup_attr"></a>
The `members` block supports:

* `id` - The ID of the backup shared member record.

* `status` - The backup sharing status.

* `created_at` - The creation time of the backup shared member.

* `updated_at` - The latest update time of the backup shared member.

* `image_id` - The ID of the image registered with the shared backup copy.

* `vault_id` - The ID of the vault where the shared backup is stored.

## Import

Share resources can be imported by their `id` or `backup_id`, e.g.

```bash
terraform import huaweicloud_cbr_backup_share.test <id>
```
