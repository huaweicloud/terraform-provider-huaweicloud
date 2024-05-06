---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_eip_associate"
description: |-
  Manages a DDS instance EIP associate resource within HuaweiCloud.
---

# huaweicloud_dds_instance_eip_associate

Manages a DDS instance EIP associate resource within HuaweiCloud.

-> **NOTE:** The shard and config nodes of a cluster instance, the read-only node of a replica set, and the hidden node
  do not support to bind the EIP.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "public_ip" {}

resource "huaweicloud_dds_instance_eip_associate" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  public_ip   = var.public_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance.
  Changing this creates a new resource.

* `node_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance node.
  Changing this creates a new resource.

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address. Changing this creates a new resource.

## Attribute Reference

* `id` - Indicates the resource ID. Format is `<instance_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `delete` - Default is 30 minutes.

## Import

DDS instance node association information can be imported using `<instance_id>/<node_id>`, e.g.

```bash
terraform import huaweicloud_dds_instance_eip_associate.test <instance_id>/<node_id>
```
