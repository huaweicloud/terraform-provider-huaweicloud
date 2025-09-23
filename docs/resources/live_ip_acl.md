---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_ip_acl"
description: |-
  Manages a Live IP address acl resource within HuaweiCloud.
---

# huaweicloud_live_ip_acl

Manages a Live IP address acl resource within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "auth_type" {}
variable "ip_auth_list" {}

resource "huaweicloud_live_ip_acl" "test" {
  domain_name  = var.domain_name
  auth_type    = var.auth_type
  ip_auth_list = var.ip_auth_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `domain_name` - (Required, String) Specifies the ingest or streaming domain name.

* `auth_type` - (Required, String) Specifies the authentication mode.
  The options are as follows:
  + **WHITE**: IP address whitelist authentication.
  + **BLACK**: IP address blacklist authentication.

* `ip_auth_list` - (Required, String) Specifies the blacklist or whitelist IP addresses. Use semicolons (;) to separate
  IP addresses, for example, **192.168.0.0;192.168.0.8**. A maximum of `100` IP addresses are allowed.
  IP network segments can be added, for example, **127.0.0.1/24**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Live IP address acl resource can be imported using `domain_name`, e.g.

```bash
$ terraform import huaweicloud_live_ip_acl.test <domain_name>
```
