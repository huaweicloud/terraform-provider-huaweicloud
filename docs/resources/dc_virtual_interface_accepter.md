---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_interface_accepter"
description: |-
  Manages a DC virtual interface accepter resource within HuaweiCloud.
---

# huaweicloud_dc_virtual_interface_accepter

Manages a DC virtual interface accepter resource within HuaweiCloud.

-> **NOTE:** Deleting a resource does not change the current receive operation.

## Example Usage

```hcl
variable "virtual_interface_id" {}

resource "huaweicloud_dc_virtual_interface_accepter" "test" {
  virtual_interface_id = var.virtual_interface_id
  action               = "ACCEPTED"
} 
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

  -> The region needs to be consistent with the region where the virtual interface created by other tenants is located.

* `virtual_interface_id` - (Required, String, ForceNew) Specifies the virtual interface ID created by other tenants.

  Changing this parameter will create a new resource.

* `action` - (Required, String, ForceNew) Specifies the action on virtual interfaces created by other tenants.
  Valid values are **ACCEPTED** and **REJECTED**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `virtual_interface_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
