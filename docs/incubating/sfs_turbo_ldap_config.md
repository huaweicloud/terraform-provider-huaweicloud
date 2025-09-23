---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_ldap_config"
description: |-
  Use this resource to manage the LDAP configuration of the SFS Turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_ldap_config

Use this resource to manage the LDAP configuration of the SFS Turbo within HuaweiCloud.

## Example Usage

```hcl
variable "share_id" {}
variable "url" {}
variable "base_dn" {}
variable "user_dn" {}
variable "password" {}

resource "huaweicloud_sfs_turbo_ldap_config" "test" {
  share_id = var.share_id
  url      = var.url
  base_dn  = var.base_dn
  user_dn  = var.user_dn
  password = var.password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, NonUpdatable) Specifies the ID of the SFS Turbo file system.

* `url` - (Required, String) Specifies the URL of the LDAP server.
  The format is `ldap://{ip_address}:{port_number}` or `ldaps://{ip_address}:{port_number}`,
  for example, **ldap://192.168.xx.xx:60000**.

* `base_dn` - (Required, String) Specifies the base distinguished name (DN) for LDAP searches.

* `user_dn` - (Optional, String) Specifies the bind DN used to authenticate to the LDAP server.

* `password` - (Optional, String) Specifies the password for the bind DN. This field is sensitive and will not be
  displayed in the state.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC that the specified LDAP server can connect to.
  This parameter is only required when the SFS Turbo file system is used across VPCs.

* `backup_url` - (Optional, String) Specifies the URL of the standby LDAP server.
  The format is `ldap://{ip_address}:{port_number}` or `ldaps://{ip_address}:{port_number}`,
  for example, **ldap://192.168.xx.xx:60000**.

* `schema` - (Optional, String) Specifies the LDAP schema. If not specified, **RFC2307** will be used.

* `search_timeout` - (Optional, Int) Specifies the LDAP search timeout interval, in seconds.
  If not specified, `3` seconds will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the `share_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SFS Turbo LDAP configuration can be imported using the `share_id`, e.g.

```bash
$ terraform import huaweicloud_sfs_turbo_ldap_config.test <share_id>
```

```
Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the imported state. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_sfs_turbo_ldap_config" "test" {
    ...
  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
