---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_authorization"
description: ""
---

# huaweicloud_cc_authorization

Manages a cross-account authorization resource within HuaweiCloud.

## Example Usage

```hcl
variable instance_id {}
variable cloud_connection_domain_id {}
variable cloud_connection_id {}

resource "huaweicloud_cc_authorization" "test" {
   name                       = "demo"
   instance_type              = "vpc"
   instance_id                = var.instance_id
   cloud_connection_domain_id = var.cloud_connection_domain_id
   cloud_connection_id        = var.cloud_connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Optional, String) The name of the cross-account authorization.

* `instance_type` - (Required, String, ForceNew) The instance type.
  The options are as follows:
    + **vpc**: VPC

  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) The instance ID.

  Changing this parameter will create a new resource.

* `cloud_connection_domain_id` - (Required, String, ForceNew) The peer account ID that you want to authorize.

  Changing this parameter will create a new resource.

* `cloud_connection_id` - (Required, String, ForceNew) Peer cloud connection ID.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the cross-account authorization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The cross-account authorization can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cc_authorization.test 0ce123456a00f2591fabc00385ff1234
```
