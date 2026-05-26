---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_driver_delete"
description: |-
  Manages a resource to delete DRS driver files within HuaweiCloud.
---

# huaweicloud_drs_driver_delete

Manages a resource to delete DRS driver files within HuaweiCloud.

-> This resource is a one-time action resource used to delete DRS driver files. Deleting this resource
   will not restore the deleted driver files or undo the delete action, but will only remove the resource information from
   the tf state file.

## Example Usage

```hcl
variable "driver_names" { 
  type = list(string)
}

resource "huaweicloud_drs_driver_delete" "test" { 
  driver_type  = "db2"
  driver_names = var.driver_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `driver_type` - (Required, String, NonUpdatable) Specifies the type of the driver file to be deleted. Valid values are:
  + **db2**: DB2 for LUW
  + **informix**: Informix

* `driver_names` - (Required, List, NonUpdatable) Specifies the list of JDBC driver file names. The list contains `1` to
  `20` files. The value of `driver_name` contains `5` to `64` characters and ends with **.jar**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
