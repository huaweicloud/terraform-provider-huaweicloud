---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_relations_details"
description: |-
  Use this data source to get the list of RMS resource relations details.
---

# huaweicloud_rms_resource_relations_details

Use this data source to get the list of RMS resource relations details.

## Example Usage

```hcl
variable "resource_id" {}
variable "direction" {}

data "huaweicloud_rms_resource_relations_details" "test" {
  resource_id = var.resource_id
  direction   = var.direction
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Required, String) Specifies the resource ID.

* `direction` - (Required, String) Specifies the direction of a resource relationship.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `relations` - The list of resource relationships.

  The [relations](#relations_struct) structure is documented below.

<a name="relations_struct"></a>
The `relations` block supports:

* `relation_type` - The relationship type.

* `from_resource_type` - The type of the source resource.

* `to_resource_type` - The type of the destination resource.

* `from_resource_id` - The source resource ID.

* `to_resource_id` - The destination resource ID.
