---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_update_security_levels_sort"
description: |-
  Manages a resource to adjust the display order of DSC sensitive data classification levels within HuaweiCloud.
---

# huaweicloud_dsc_update_security_levels_sort

Manages a resource to adjust the display order of DSC sensitive data classification levels within HuaweiCloud.

-> This resource is a one-time action resource used to adjust the display order of DSC security levels. Deleting this
  resource will not revert the sort adjustment, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "level_id" {}
variable "target_level_id" {}

resource "huaweicloud_dsc_update_security_levels_sort" "test" {
  level_id        = var.level_id
  target_level_id = var.target_level_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `level_id` - (Required, String, NonUpdatable) Specifies the unique ID of the security level to be adjusted.

* `target_level_id` - (Optional, String, NonUpdatable) Specifies the ID of the target reference security level,
  used to locate the sort position.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message, which describes the operation result or status information.

* `status` - The returned status, indicating whether the operation is successful.
