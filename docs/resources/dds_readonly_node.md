---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_readonly_node"
description: |-
  Manages a DDS readonly node resource within HuaweiCloud.
---

# huaweicloud_dds_readonly_node

Manages a DDS readonly node resource within HuaweiCloud.

-> Before use this resource, you need to pay attention to the following:
  <br/>1. The resource only supports replica set DDS instance to add readonly node.
  <br/>2. A maximum of `5` readonly node can be added to a replica set DDS instance.

## Example Usage

```hcl
variable "instance_id" {}
variable "spec_code" {}

resource "huaweicloud_dds_readonly_node" "test"{
  instance_id = var.instance_id
  spec_code   = var.spec_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DDS instance.

* `spec_code` - (Required, String) Specifies the readonly node specification code.

* `delay` - (Optional, Int, NonUpdatable) Specifies the synchronization delay time.
  The valid value ranges from `0` to `1,200`, default is `0`, unit is ms.

* `size` - (Optional, String) Specifies the requested disk capacity.
  The value must be an integer multiple of `10` and greater than the current storage space.

* `private_ip` - (Optional, String) Specifies the readonly node private IP address.
  Currently, only the IPv4 address is supported. An in-use IP address cannot be used as the private IP address of
  the DDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - The readonly node name.

* `status` - The readonly node status.

* `role` - The node role.

* `availability_zone` - The availability zone to which the readonly node belongs.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 30 minutes.

## Import

The resource can be imported using the `instance_id` and the `id` separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dds_readonly_node.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `delay`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dds_readonly_node" "test" {
  ...

  lifecycle {
    ignore_changes = [
      delay,
    ]
  }
}
```
