---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_enable_trusted_service"
description: |-
  Manages a resource to enable the organization trusted service within HuaweiCloud.
---

# huaweicloud_dsc_multi_enable_trusted_service

Manages a resource to enable the organization trusted service within HuaweiCloud.

-> This resource is a one-time action resource used to enable the organization trusted service. Deleting this
  resource will not disable the trusted service or undo the enable action, but will only remove the resource
  information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_dsc_multi_enable_trusted_service" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
