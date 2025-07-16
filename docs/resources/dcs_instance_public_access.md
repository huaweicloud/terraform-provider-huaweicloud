---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_public_access"
description: |-
  Manages a DCS instance public access resource within HuaweiCloud.
---

# huaweicloud_dcs_instance_public_access

Manages a DCS instance public access resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "elb_id" {}

resource "huaweicloud_dcs_instance_public_access" "test"{
  instance_id = var.instance_id
  elb_id      = var.elb_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `publicip_id` - (Optional, String, NonUpdatable) Specifies the ID of the public IP address. This parameter is mandatory
  when **Redis 3.0** is used.

* `enable_ssl` - (Optional, Bool, NonUpdatable) Specifies whether to enable SSL. This parameter has a value only when SSL
  is enabled. This parameter is mandatory for **Redis 3.0**.

* `elb_id` - (Optional, String, NonUpdatable) Specifies the ID of the load balancer bound for public access. This
  parameter is mandatory when **Redis 4.0** or later is used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `instance_id`.

* `eip_id` - Indicates the ID of the EIP.

* `eip_address` - Indicates the address of the EIP.

* `elb_listeners` - Indicates the list of the ELB listeners.
  The [elb_listeners](#elb_listeners_structs) structure is documented below.

<a name="elb_listeners_structs"></a>
The `elb_listeners` block supports:

* `id` - Indicates the ID of the listener.

* `port` - Indicates the port of the listener.

* `name` - Indicates the name of the listener.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The DCS instance public access can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_dcs_instance_public_access.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `publicip_id` and `enable_ssl`. It is
generally recommended running `terraform plan` after importing a resource. You can then decide if changes should be
applied to the resource, or the resource definition should be updated to align with the resource. Also, you can ignore
changes as below.

```bash
resource "huaweicloud_dcs_instance_public_access" "test" {
    ...

  lifecycle {
    ignore_changes = [
      publicip_id, enable_ssl,
    ]
  }
}
```
