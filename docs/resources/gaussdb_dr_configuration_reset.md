---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_configuration_reset"
description: |-
  Use this resource to reset the DR configuration for a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_configuration_reset

Use this resource to reset the DR configuration for a GaussDB instance within HuaweiCloud.

-> This resource is a one-time action resource for resetting the GaussDB DR configuration. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_dr_configuration_reset" "test" {
  instance_id        = var.instance_id
  opposite_data_cidr = "192.168.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `opposite_data_cidr` - (Optional, String, NonUpdatable) Specifies the subnet IP CIDR whitelist configuration of the
  remote instance. Multiple CIDRs are separated by commas. When set to an empty string, it represents clearing the
  whitelist configuration.

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.
