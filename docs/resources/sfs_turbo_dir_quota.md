---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_dir_quota"
description: |-
  Manages a SFS Turbo directory quota resource within HuaweiCloud.
---

# huaweicloud_sfs_turbo_dir_quota

Manages a SFS Turbo directory quota resource within HuaweiCloud.

## Example Usage

```hcl
variable "share_id" {}
variable "path" {}

resource "huaweicloud_sfs_turbo_dir_quota" "test" {
  path      = var.path
  share_id  = var.share_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo directory resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo directory resource.

* `share_id` - (Required, String, ForceNew) Specifies the SFS Turbo ID. Changing this will create a new resource.

* `path` - (Required, String, ForceNew) Specifies the valid full path of existing SFS Turbo directory. The parameter
  starts with "/", otherwise the parameter is illegal. Changing this will create a new resource.

* `capacity` - (Optional, Int) Specifies the size of the directory. The default value is `0`, this value means that
  there is no quota assigned to this directory. The value can not exceed the remaining available capacity of the
  file system. The unit is `MB`.

* `inode` - (Optional, Int) Specifies the maximum number of inodes allowed in the directory. The default value is `0`,
  this value means that there is no quota assigned to this directory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `used_capacity` - The size of the used directory. The unit is `MB`. This parameter is returned only for SFS Turbo
  HPC file systems.

* `used_inode` - The number of used inodes in the directory. This parameter is returned only for SFS Turbo
  HPC file systems.
