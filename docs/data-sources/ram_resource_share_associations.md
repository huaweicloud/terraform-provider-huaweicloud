---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_share_associations"
description: |-
  Use this data source to get the list of principals and resources associated with the RAM shared resources.
---

# huaweicloud_ram_resource_share_associations

Use this data source to get the list of principals and resources associated with the RAM shared resources.

## Example Usage

```hcl
variable "association_type" {}

data "huaweicloud_ram_resource_share_associations" "test" {
  association_type = var.association_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `association_type` - (Required, String) Specifies the association type. Valid values are **principal** and **resource**.

* `principal` - (Optional, String) Specifies the principal associated with the resource share.

* `resource_urn` - (Optional, String) Specifies the URN of the resource associated with the resource share.

* `resource_share_ids` - (Optional, List) Specifies the list of resource share IDs.

* `resource_ids` - (Optional, List) Specifies the list of resource IDs.

* `status` - (Optional, String) Specifies the status of the association.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `associations` - The list of association details.

  The [associations](#associations_struct) structure is documented below.

<a name="associations_struct"></a>
The `associations` block supports:

* `associated_entity` - The associated entity. It can be the resource URN, account ID, URN of the root OU, or URN of
  another OU.

* `created_at` - The time when the association was created.

* `status` - The status of the association.

* `status_message` - The description of the status to the association.

* `association_type` - The entity type in the association.

* `updated_at` - The time when the association was last updated.

* `external` - Whether the principle is in the same organization as the resource owner.

* `resource_share_id` - ID of the resource share.

* `resource_share_name` - Name of the resource share.
