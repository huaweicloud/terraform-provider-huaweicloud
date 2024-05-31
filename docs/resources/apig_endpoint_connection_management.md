---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_endpoint_connection_management"
description: |-
  Use this resource to accept or reject a VPC endpoint connection under specified instance within HuaweiCloud.
---

# huaweicloud_apig_endpoint_connection_management

Use this resource to accept or reject a VPC endpoint connection under specified instance within HuaweiCloud.

-> Destroying resources does not change the current action of the endpoint connection.

## Example Usage

```hcl
variable instance_id {}
variable endpoint_id {}

resource "huaweicloud_apig_endpoint_connection_management" "test" {
  instance_id = var.instance_id
  action      = "receive"
  endpoint_id = var.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the dedicated instance ID to which the endpoint connection belongs.
  Changing this creates a new resource.

* `endpoint_id` - (Required, String, ForceNew) Specifies the ID of the endpoint connection.
  Changing this creates a new resource.
  
* `action` - (Required, String) Specifies the operation type endpoint connection.
  The valid values are as follows:
  + **receive**
  + **reject**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The current ststus of the endpoint connection.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
