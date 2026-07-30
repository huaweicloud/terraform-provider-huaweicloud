---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_metadata_tags"
description: |-
  Manages DSC metadata tags resource within HuaweiCloud.
---

# huaweicloud_dsc_metadata_tags

Manages DSC metadata tags resource within HuaweiCloud.

## Example Usage

```hcl
variable "tag_names" {
  type = list(string)
}

resource "huaweicloud_dsc_metadata_tags" "test" {
  names = var.tag_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `names` - (Required, List) Specifies the list of tag names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
