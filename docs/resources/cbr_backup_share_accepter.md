---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_backup_share_accepter"
description: ""
---

# huaweicloud_cbr_backup_share_accepter

Using this resource to accept a shared backup within HuaweiCloud.

## Example Usage

```hcl
variable "backup_id" {}
variable "stored_vault_id" {}

resource "huaweicloud_cbr_backup_share_accepter" "test" {
  backup_id = var.backup_id
  vault_id  = var.stored_vault_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the backup will be stored.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `backup_id` - (Required, String, ForceNew) Specifies the ID of the shared source backup.  
  Changing this will create a new resource.

* `vault_id` - (Required, String, ForceNew) Specifies the ID of the vault which the backup will be stored.  
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also backup ID) in UUID format.

* `source_project_id` - The ID of the project to which the source backup belongs.

## Import

Resources can be imported by their `id` or `backup_id`, e.g.

```bash
$ terraform import huaweicloud_cbr_backup_share_accepter.test <id>
```
