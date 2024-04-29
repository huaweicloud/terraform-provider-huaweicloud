---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_proxy"
description: ""
---

# huaweicloud_gaussdb_mysql_proxy

GaussDB mysql proxy management within HuaweiCould.

## Example Usage

### create a proxy

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_proxy" "proxy_1" {
  instance_id = var.instance_id
  flavor      = "gaussdb.proxy.xlarge.arm.2"
  node_num    = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql proxy resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID of the proxy.
  Changing this parameter will create a new resource.

* `flavor` - (Required, String, ForceNew) Specifies the flavor of the proxy.
  Changing this parameter will create a new resource.

* `node_num` - (Required, Int) Specifies the node count of the proxy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID in UUID format.
* `address` - Indicates the address of the proxy.
* `port` - Indicates the port of the proxy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

GaussDB instance can be imported using the instance `id`, e.g.

```
$ terraform import huaweicloud_gaussdb_mysql_proxy.proxy_1 ee678f40-ce8e-4d0c-8221-38dead426f06
```
