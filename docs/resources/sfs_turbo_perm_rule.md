---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_perm_rule"
description: |-
  Manages a SFS Turbo permission rule resource within HuaweiCloud.
---

# huaweicloud_sfs_turbo_perm_rule

Manages a SFS Turbo permission rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "share_id" {}

resource "huaweicloud_sfs_turbo_perm_rule" "test" {
  share_id  = var.share_id
  ip_cidr   = "192.168.0.0/16"
  rw_type   = "rw"
  user_type = "no_root_squash"
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo permission rule resource.
  If omitted, the provider-level region will be used. Changing this creates a new SFS Turbo permission rule resource.

* `share_id` - (Required, String, ForceNew) Specifies the SFS Turbo ID. Changing this will create a new resource.

* `ip_cidr` - (Required, String, ForceNew) Specifies the IP address or IP address range of the object to be authorized.
  Changing this will create a new SFS Turbo permission rule resource.

* `rw_type` - (Required, String) Specifies the read/write permission of the object to be authorized.
  The value can be **rw** (read and write permission) or **ro** (read only permission).

* `user_type` - (Required, String) Specifies the file system access permission granted to the user of the object to be
  authorized. The value can be **no_root_squash**, **root_squash** or **all_squash**.
  + **no_root_squash** allows the root user on the client to access the file system as **root**.
  + **root_squash** allows the root user on the client to access the file system as **nfsnobody**.
  + **all_squash** allows any user on the client to access the file system as **nfsnobody**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
