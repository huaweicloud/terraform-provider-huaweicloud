---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_image_shared_members"
description: |-
  Use this data source to get the list of IMS image shared members within HuaweiCloud.
---

# huaweicloud_ims_image_shared_members

Use this data source to get the list of IMS image shared members within HuaweiCloud.

## Example Usage

```hcl
variable "image_id"{}

data "huaweicloud_ims_image_shared_members" "test" {
  image_id = var.image_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `image_id` - (Required, String) Specifies the image ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `members` - Indicates the members.

  The [members](#members_struct) structure is documented below.

* `schema` - Indicates the schema.

<a name="members_struct"></a>
The `members` block supports:

* `member_id` - Indicates the member ID.

* `member_type` - Indicates the member type.

* `urn` - Indicates the shared organization urn.

* `status` - Indicates the shared status.

* `created_at` - Indicates the shared time.

* `updated_at` - Indicates the update time.

* `image_id` - Indicates the image ID.

* `schema` - Indicates the schema.
