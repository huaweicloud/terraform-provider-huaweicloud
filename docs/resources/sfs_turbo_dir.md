---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_dir"
description: |-
  Manages a SFS Turbo directory resource within HuaweiCloud.
---

# huaweicloud_sfs_turbo_dir

Manages a SFS Turbo directory resource within HuaweiCloud.

## Example Usage

```hcl
variable "share_id" {}

resource "huaweicloud_sfs_turbo_dir" "test" {
  path     = "/tmp01"
  share_id = var.share_id
  mode     = 777
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo directory resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo directory resource.

* `share_id` - (Required, String, ForceNew) Specifies the SFS Turbo ID. Changing this will create a new resource.

* `path` - (Required, String, ForceNew) Specifies the valid full path of SFS Turbo directory. The parameter
  starts with "/", otherwise the parameter is illegal. Changing this will create a new resource.

* `mode` - (Optional, Int, ForceNew) Specifies the file directory permissions. The valid value ranges from `0` to `777`.
  `0` means no permission, and `777` means the highest authority(read, write and execute). Please refer
  to [document](https://en.wikipedia.org/wiki/Chmod#Numerical_permissions). This field will be set to `0` when
  this field over `777`.

* `gid` - (Optional, Int, ForceNew) Specifies the group ID of the file directory. The minimum value is `0`,
  the value represents the ID of the super user `root`.

* `uid` - (Optional, Int, ForceNew) Specifies the user ID of the file directory. The minimum value is `0`,
  the value represents the ID of the group where the super user `root` belongs.
