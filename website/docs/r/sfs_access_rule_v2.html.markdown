---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_access_rule_v2"
sidebar_current: "docs-huaweicloud-resource-sfs-access-rule-v2"
description: |-
 Provides an access rule resource of Scalable File Resource (SFS).
---

# huaweicloud_sfs_access_rule_v2

Provides an access rule resource of Scalable File Resource (SFS).

## Example Usage

### Usage in VPC authorization scenarios
```hcl
variable "share_name" { }
variable "vpc_id" { }

resource "huaweicloud_sfs_file_system_v2" "share-file" {
  name        = var.share_name
  size        = 100
  share_proto = "NFS"
}

resource "huaweicloud_sfs_access_rule_v2" "rule_1" {
  sfs_id    = huaweicloud_sfs_file_system_v2.share-file.id
  access_to = var.vpc_id
}
```

### Usage in IP address authorization scenario
```hcl
variable "share_name" { }
variable "vpc_id" { }

resource "huaweicloud_sfs_file_system_v2" "share-file" {
  name        = var.share_name
  size        = 100
  share_proto = "NFS"
}

resource "huaweicloud_sfs_access_rule_v2" "rule_1" {
  sfs_id    = huaweicloud_sfs_file_system_v2.share-file.id
  access_to = join("#", [var.vpc_id, "192.168.10.0/24", "0", "no_all_squash,no_root_squash"])
}
```

## Argument Reference
The following arguments are supported:

* `sfs_id` - (Required) Specifies the UUID of the shared file system. Changing this will create a new access rule.

* `access_level` - (Optional) Specifies the access level of the shared file system. Possible values are *ro* (read-only)
    and *rw* (read-write). The default value is *rw* (read/write). Changing this will create a new access rule.

* `access_type` - (Optional) Specifies the type of the share access rule. The default value is *cert*.
    Changing this will create a new access rule.

* `access_to` - (Required) Specifies the value that defines the access rule. The value contains 1 to 255 characters.
    Changing this will create a new access rule. The value varies according to the scenario:
    - Set the VPC ID in VPC authorization scenarios.
    - Set this parameter in IP address authorization scenario.

        For an NFS shared file system, the value in the format of *VPC_ID#IP_address#priority#user_permission*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#100#all_squash,root_squash.

        For a CIFS shared file system, the value in the format of *VPC_ID#IP_address#priority*.
        For example, 0157b53f-4974-4e80-91c9-098532bcaf00#2.2.2.2/16#0.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the share access rule.

* `status` - The status of the share access rule.

## Import

SFS access rule can be imported by specifying the SFS ID and access rule ID separated by a slash, e.g.:

```
$ terraform import huaweicloud_sfs_access_rule_v2 <sfs_id>/<rule_id>
```