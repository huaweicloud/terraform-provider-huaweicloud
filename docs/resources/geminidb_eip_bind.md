---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_eip_bind"
description: |-
  Manages a resource to bind the EIP to the GeminiDB instance node within HuaweiCloud.
---

# huaweicloud_geminidb_eip_bind

Manages a resource to bind the EIP to the GeminiDB instance node within HuaweiCloud.

## Example Usage

```hcl
var "instance_id" {}
var "node_id" {}
var "eip_ip_address" {}
var "eip_id" {}

resource "huaweicloud_geminidb_eip_bind" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  public_ip    = var.eip_ip_address
  public_ip_id = eip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the node.

* `public_ip` - (Required, String, NonUpdatable) Specifies the EIP IP address.

* `public_ip_id` - (Required, String, NonUpdatable) Specifies the EIP ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the node ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The resource can be imported using the `instance_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_geminidb_eip_bind.test <instance_id>/<id>
```
