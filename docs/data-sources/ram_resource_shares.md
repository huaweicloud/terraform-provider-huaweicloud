---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_shares"
description: |-
  Use this data source to get a list of resource shares that you have created or shared with you.
---

# huaweicloud_ram_resource_shares

Use this data source to get a list of resource shares that you have created or shared with you.

## Example Usage

```hcl
variable "resource_owner" {}

data "huaweicloud_ram_resource_shares" "test" {
  resource_owner = var.resource_owner
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_owner` - (Required, String) Specifies the owner type of resource sharing instance.
  Valid values are **self** and **other-accounts**.

* `resource_share_ids` - (Optional, List) Specifies the list of resource share IDs.

* `status` - (Optional, String) Specifies the status of the resource share.

* `tag_filters` - (Optional, List) Specifies the tags attached to the resource share.

  The [tag_filters](#tag_filters_struct) structure is documented below.

* `name` - (Optional, String) Specifies the name of the resource share.

* `permission_id` - (Optional, String) Specifies the permission ID.

<a name="tag_filters_struct"></a>
The `tag_filters` block supports:

* `key` - (Required, String) Specifies the identifier or name of the tag key.

* `values` - (Optional, List) Specifies the list of values for the tag key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_shares` - The list of details about resource shares.

  The [resource_shares](#resource_shares_struct) structure is documented below.

<a name="resource_shares_struct"></a>
The `resource_shares` block supports:

* `name` - The name of the resource share.

* `allow_external_principals` - Whether resources can be shared with any accounts outside the organization.

* `tags` - The list of tags attached to the resource share.

  The [tags](#resource_shares_tags_struct) structure is documented below.

* `created_at` - The time when the resource share was created.

* `updated_at` - The time when the resource share was last updated.

* `id` - The ID of the resource share.

* `owning_account_id` - The ID of the resource owner in a resource share.

* `status` - The Status of the resource share.

* `description` - The description of the resource share.

<a name="resource_shares_tags_struct"></a>
The `tags` block supports:

* `key` - Identifier or name of the tag key.

* `value` - Tag value.
