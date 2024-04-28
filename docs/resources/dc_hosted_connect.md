---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_hosted_connect"
description: ""
---

# huaweicloud_dc_hosted_connect

Manages a hosted connection resource within HuaweiCloud.

-> The creator **must** have the partner qualification and have an operations connection.

## Example Usage

```hcl
variable resource_tenant_id {}
variable hosting_id {}

resource "huaweicloud_dc_hosted_connect" "test" {
  name               = "demo"
  description        = "This is a demo"
  resource_tenant_id = var.resource_tenant_id
  hosting_id         = var.hosting_id
  vlan               = 441
  bandwidth          = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Optional, String) The name of the hosted connect.

* `bandwidth` - (Required, Int) The bandwidth size of the hosted connect in Mbit/s.

* `hosting_id` - (Required, String, ForceNew) The ID of the operations connection on which the hosted connect is created.

  Changing this parameter will create a new resource.

* `vlan` - (Required, Int, ForceNew) The VLAN allocated to the hosted connect.

  Changing this parameter will create a new resource.

* `resource_tenant_id` - (Required, String, ForceNew) The tenant ID for whom a hosted connect is to be created.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the hosted connect.

* `peer_location` - (Optional, String) The location of the on-premises facility at the other end of the connection.  
  Specific to the street or data center name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the hosted connect.  
  The options are as follows:
  + **BUILD**: The hosted connect has been created.
  + **ACTIVE**: The associated virtual gateway is normal.
  + **DOWN**: The port used by the hosted connect is down, indicating that there may be line faults.
  + **ERROR**: The associated virtual gateway is abnormal.
  + **PENDING_DELETE**: The hosted connect is being deleted.
  + **PENDING_UPDATE**: The hosted connect is being updated.
  + **PENDING_CREATE**: The hosted connect is being created.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The hosted connect can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dc_hosted_connect.test ac0fe389-02f5-4463-9647-58bbb3d21fed
```
