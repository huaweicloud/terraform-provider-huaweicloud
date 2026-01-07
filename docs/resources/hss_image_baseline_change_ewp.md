---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_baseline_change_ewp"
description: |-
  Using this resource to update the extended-weak-password configuration of the HSS image baseline within HuaweiCloud.
---

# huaweicloud_hss_image_baseline_change_ewp

Using this resource to update the extended weak password configuration of the HSS image baseline within HuaweiCloud.

-> This is a one-time action resource used to update the extended-weak-password configuration of the HSS image baseline.
  Deleting this resource will not clear the corresponding request record, but will only remove the resource information
  from the tf state file.

## Example Usage

```hcl
variable "extended_weak_password" {
  type = list(string)
}

resource "huaweicloud_hss_image_baseline_change_ewp" "test" {
  extended_weak_password = var.extended_weak_password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `extended_weak_password` - (Optional, List, NonUpdatable) Specifies the extended-weak-password string list.  

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
