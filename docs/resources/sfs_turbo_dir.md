---
subcategory: "Scalable File Service (SFS)"
---

# huaweicloud_sfs_turbo_dir

Provides an Shared File System (SFS) Turbo directory resource.

## Example Usage

### Create a STANDARD Shared File System (SFS) Turbo Directory

```hcl
variable "share_id" {}

resource "huaweicloud_sfs_turbo_dir" "test" {
  path     = "/tmp01"
  share_id = var.share_id
  mude     = 777
  gid      = 100
  uid      = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo directory resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo directory resource.

* `share_id` - (Required, String, ForceNew) Specifies the created file system ID. Changing this will create a new resource.

* `path` - (Required, String, ForceNew) Specifies the valid full path of an existing directory. The parameter starts with "/",
  otherwise the parameter is illegal. Changing this will create a new resource.

* `mude` - (Optional, Int, ForceNew) Specifies the file directory permissions. The minimum value is 0.

* `gid` - (Optional, Int, ForceNew) Specifies the group ID of the file directory. The minimum value is 0.

* `uid` - (Optional, Int, ForceNew) Specifies the user ID of the file directory. The minimum value is 0.