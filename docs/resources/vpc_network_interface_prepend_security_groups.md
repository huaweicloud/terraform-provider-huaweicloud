---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eni_sg_attachment"
description: ""
---

# huaweicloud_vpc_eni_sg_attachment

Manages the prepend order of security groups for an existing network interface (normal or supplemental).

On create/update, this resource prepends `prepend_security_group_ids` before the current security group order and de-duplicates IDs.
On delete, it restores the original security group bindings captured at create time.

## Example Usage

```hcl
resource "huaweicloud_vpc_eni_sg_attachment" "test" {
  network_interface_id   = "7f49c11a-xxxx-xxxx-xxxx-2a19127d1a5c"
  network_interface_type = "normal"

  prepend_security_group_ids = [
    "sg-4",
    "sg-3",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query and update the resource.
  If omitted, the provider-level region will be used.

* `network_interface_id` - (Required, String, ForceNew) Specifies the target network interface ID.

* `network_interface_type` - (Optional, String, ForceNew) Specifies the target network interface type.
  The value can be:
  + **normal**: normal elastic network interface.
  + **supplemental**: supplemental network interface.
  Defaults to **normal**.

* `prepend_security_group_ids` - (Required, List) Specifies the security group IDs to prepend.
  Existing IDs in current bindings will be moved to the front according to this order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<network_interface_type>/<network_interface_id>`.

* `original_security_group_ids` - The original security group bindings captured when this resource is created.

* `effective_security_group_ids` - The current effective security group bindings after prepend operation.

## Import

The resource can be imported using `<network_interface_type>/<network_interface_id>`, e.g.

```bash
$ terraform import huaweicloud_vpc_eni_sg_attachment.test normal/7f49c11a-xxxx-xxxx-xxxx-2a19127d1a5c
```
