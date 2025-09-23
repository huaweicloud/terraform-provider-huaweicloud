---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_file_system"
description: ""
---

# huaweicloud_sfs_file_system

!> **WARNING:** It has been deprecated.

Provides a Shared File System (SFS) resource within HuaweiCloud.

## Example Usage

### Basic example

```hcl
variable "share_name" {}
variable "share_description" {}
variable "vpc_id" {}

resource "huaweicloud_sfs_file_system" "share-file" {
  name         = var.share_name
  size         = 100
  share_proto  = "NFS"
  access_level = "rw"
  access_to    = var.vpc_id
  description  = var.share_description

  tags = {
    key = "value"
  }
}
```

### SFS with data encryption

```hcl
variable "share_name" {}
variable "share_description" {}
variable "vpc_id" {}

respurce "huaweicloud_kms_key" mykey {
  key_alias    = "kms_sfs"
  pending_days = "7"
}

resource "huaweicloud_sfs_file_system" "share-file" {
  name         = var.share_name
  size         = 100
  share_proto  = "NFS"
  access_level = "rw"
  access_to    = var.vpc_id
  description  = var.share_description

  metadata = {
    "#sfs_crypt_key_id"    = huaweicloud_kms_key.mykey.id
    "#sfs_crypt_domain_id" = huaweicloud_kms_key.mykey.domain_id
    "#sfs_crypt_alias"     = huaweicloud_kms_key.mykey.key_alias
  }
  tags     = {
    function = "encryption"
  }
}
```

### SFS with Auto Capacity Expansion

-> This feature is only supported in specific regions.

```hcl
variable "share_name" {}

resource "huaweicloud_sfs_file_system" "share-file" {
  name        = var.share_name
  size        = 100
  share_proto = "NFS"
  description = "auto capacity expansion"

  metadata = {
    "#sfs_quota_type" = "soft"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the sfs resource. If omitted, the provider-level
  region will be used. Changing this creates a new sfs resource.

* `size` - (Required, Int) Specifies the size (GB) of the shared file system.

* `share_proto` - (Optional, String) Specifies the file system sharing protocol.
  The valid value can be **NFS** (for Linux OS) or **CIFS** (for Windows OS).

* `name` - (Optional, String) Specifies the name of the shared file system, which contains 0 to 255 characters and
  can contain only letters, digits, hyphens (-), and underscores (_).

* `description` - (Optional, String) Specifies the description of the shared file system, which contains 0 to 255
  characters and can contain only letters, digits, hyphens (-), and underscores (_).

* `is_public` - (Optional, Bool, ForceNew) Specifies whether a file system can be publicly seen.
  The default value is false.

* `metadata` - (Optional, Map, ForceNew) Specifies the metadata information used to create the shared file system. The
  supported metadata keys are "#sfs_crypt_key_id", "#sfs_crypt_domain_id" and "#sfs_crypt_alias", and the keys should be
  exist at the same time to enable the data encryption function. Changing this will create a new resource.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone name. Changing this parameter will
  create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the shared file system. Changing
  this creates a new resource.

* `tags` - (Optional, Map) The key/value pairs to associate with the shared file system.

* `access_level` - (Optional, String) Specifies the access level of the shared file system. Possible values are *ro* (
  read-only)
  and *rw* (read-write). The default value is *rw* (read/write). Changing this will create a new access rule.

* `access_type` - (Optional, String) Specifies the type of the share access rule. The default value is *cert*. Changing
  this will create a new access rule.

* `access_to` - (Optional, String) Specifies the value that defines the access rule. The value contains 1 to 255
  characters. Changing this will create a new access rule. The value varies according to the scenario:
  + Set the VPC ID in VPC authorization scenarios.
  + Set this parameter in IP address authorization scenario:
      - For an NFS shared file system, the value in the format of *VPC_ID#IP_address#priority#user_permission*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#100#all_squash,root_squash.
      - For a CIFS shared file system, the value in the format of *VPC_ID#IP_address#priority*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#0.

-> **NOTE:** If you want to create more access rules, please using
  [huaweicloud_sfs_access_rule](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/resources/sfs_access_rule).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the shared file system.

* `status` - The status of the shared file system.

* `export_location` - The address for accessing the shared file system.

* `share_access_id` - The UUID of the share access rule.

* `access_rule_status` - The status of the share access rule.

* `access_rules` - All access rules of the shared file system. The object includes the following:
  + `access_rule_id` - The UUID of the share access rule.
  + `access_level` - The access level of the shared file system
  + `access_type` - The type of the share access rule.
  + `access_to` - The value that defines the access rule.
  + `status` - The status of the share access rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

SFS can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_sfs_file_system 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```

**NOTE:** The `access_to`, `access_type` and `access_level` will not be imported. Please importing them by
  [huaweicloud_sfs_access_rule](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/resources/sfs_access_rule).
