---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_internal_ip_modify"
description: |-
  Manages a DDS instance internal IP modify resource within HuaweiCloud.
---

# huaweicloud_dds_instance_internal_ip_modify

Manages a DDS instance internal IP modify resource within HuaweiCloud.

-> **NOTE:** Deleting instance internal IP modify is not supported. If you destroy a resource of instance internal IP
  modify, it is only removed from the state, but still remains in the cloud. And the instance doesn't return to the
  state before modifying.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "new_ip" {}

resource "huaweicloud_dds_instance_internal_ip_modify" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  new_ip      = var.new_ip
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

* `new_ip` - (Required, String) Specifies the new IP of the DDS instance node.

-> By default, only the internal IP of the mongos nodes and replica set instances nodes can be modified.

## Attribute Reference

* `id` - Indicates the resource ID. Format is `<instance_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 30 minutes.

## Import

DDS instance node can be imported using `<instance_id>/<node_id>`, e.g.

```bash
terraform import huaweicloud_dds_instance_internal_ip_modify.test <instance_id>/<node_id>
```
