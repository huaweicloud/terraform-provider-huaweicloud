---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_proxy_eip_associate"
description: |-
  Manages TaurusDB Proxy instance EIP associate resource within HuaweiCloud.
---

# huaweicloud_taurusdb_proxy_eip_associate

Manages TaurusDB Proxy instance EIP associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
vaiable "proxy_id" {}
variable "public_ip" {}
variable "public_ip_id" {}

resource "huaweicloud_taurusdb_proxy_eip_associate" "test"{
  instance_id  = var.instance_id
  proxy_id     = var.proxy_id
  public_ip    = var.public_ip
  public_ip_id = var.public_ip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the ID of a TaurusDB instance.

* `proxy_id` - (Required, String, NoneUpdatable) Specifies the ID of a TaurusDB proxy instance.

* `public_ip` - (Required, String, NoneUpdatable) Specifies the EIP address to be bound.

* `public_ip_id` - (Required, String, NoneUpdatable) Specifies the ID of the EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The TaurusDB Proxy EIP associate can be imported using the `instance_id` and `proxy_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_taurusdb_proxy_eip_associate.test <instance_id>/<proxy_id>
```
