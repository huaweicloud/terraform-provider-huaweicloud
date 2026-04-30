---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_backup_download_policy"
description: |-
  Manages a DDS backup download policy resource within HuaweiCloud.
---

# huaweicloud_dds_backup_download_policy

Manages a DDS backup download policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "action" {}

resource "huaweicloud_dds_backup_download_policy" "test"{
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, String) Specifies the backup download switch.
  The valid values are as follows:
  + **open**: The backup download function is enabled, and backups can be downloaded from the Internet.
  + **close**: The backup download function is disabled, and backups cannot be downloaded from the Internet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the project ID.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dds_backup_download_policy.test <id>
```
