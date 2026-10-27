---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_relations"
description: |-
  Use this data source to get the relations of a resource within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_relations

Use this data source to get the relations of a resource within HuaweiCloud Resource Center.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_rc_resource_relations" "test" {
  resource_id = var.resource_id
}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Required, String) Specifies the ID of the resource.

* `related_resource_id` - (Optional, String) Specifies the ID of the related resource.

* `related_resource_type` - (Optional, String) Specifies the type of the related resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `relations` - The resource relations.
  The [relations](#relations_struct) structure is documented below.

<a name="relations_struct"></a>
The relations block supports:

* `resource_id` - The resource ID.

* `resource_type` - The resource type.

* `related_resource_type` - The type of the related resource.

* `related_resource_id` - The ID of the related resource.
