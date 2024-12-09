---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_tags"
description: |-
  Use this data source to get the list of RMS resource tags.
---

# huaweicloud_rms_resource_tags

Use this data source to get the list of RMS resource tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_rms_resource_tags" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid value can be **config:policyAssignments**.

* `resource_id` - (Optional, String) Specifies the resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The tag values.
