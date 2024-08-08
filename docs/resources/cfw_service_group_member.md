---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_service_group_member"
description: ""
---

# huaweicloud_cfw_service_group_member

Manages a CFW service group member resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "protocol" {}
variable "source_port" {}
variable "dest_port" {}

resource "huaweicloud_cfw_service_group_member" "test" {
  group_id    = var.group_id
  protocol    = var.protocol
  source_port = var.source_port
  dest_port   = var.dest_port
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the service group.

  Changing this parameter will create a new resource.

* `protocol` - (Required, Int, ForceNew) Specifies the protocol type.
  The valid values are:
    + **6**: indicates TCP;
    + **17**: indicates UDP;
    + **1**: indicates ICMP;
    + **58**: indicates ICMPv6;
    + **-1**: indicates any protocol.

  Changing this parameter will create a new resource.

* `source_port` - (Required, String, ForceNew) Specifies the source port.

  Changing this parameter will create a new resource.

* `dest_port` - (Required, String, ForceNew) Specifies the destination port.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the service group member description.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The service group member can be imported using service group ID and member ID, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_service_group_member.test <group_id>/<member_id>
```
