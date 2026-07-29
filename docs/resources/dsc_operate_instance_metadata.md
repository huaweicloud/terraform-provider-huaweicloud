---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_operate_instance_metadata"
description: |-
  Manages a resource to operate the metadata of a DSC data instance asset within HuaweiCloud.
---

# huaweicloud_dsc_operate_instance_metadata

Manages a resource to operate the metadata of a DSC data instance asset within HuaweiCloud.

-> This resource is a one-time action resource used to operate the metadata of a DSC data instance asset.
  Deleting this resource will not restore the metadata or undo the operation, but will only remove the resource
  information from the tf state file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dsc_operate_instance_metadata" "test" {
  instance_id = var.instance_id
  action      = "REFRESH"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the data instance asset.

* `action` - (Required, String, NonUpdatable) Specifies the operation type.
  The valid values are as follows:
  + **REFRESH**: Refresh the metadata.
  + **DELETE**: Delete the metadata.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status.
