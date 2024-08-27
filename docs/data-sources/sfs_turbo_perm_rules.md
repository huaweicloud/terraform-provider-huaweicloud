---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_perm_rules"
description: |-
  Use this datasource to get a list of permission rules.
---

# huaweicloud_sfs_turbo_perm_rules

Use this datasource to get a list of permission rules.

## Example Usage

```hcl
variable "share_id" {}

data "huaweicloud_sfs_turbo_perm_rules" "test" {
  share_id = var.share_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS Turbo file system to which the permission rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the permission rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the permission rule.

* `ip_cidr` - The IP address or IP address range of the authorized object.

* `rw_type` - The read and write permission of the authorized object.
  The value can be **rw** (read and write), **ro** (read only) or **none** (no permission).
  The default value is **rw**.

* `user_type` - The file system access permission granted to the user of the authorized object.
  The valid values are as follow:
  + **no_root_squash**: Allow the root user on the client to access the file system as root.
  + **root_squasg**: Allow the root user on the client to access the file system as anonymous (nfsnobody).
  + **all_squash**: Allow any user on the client to access the file system as nfsnobody.
  
  The default value is **all_squash**.
