---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_organizations"
description: ""
---

# huaweicloud_swr_organizations

Use this data source to get the list of SWR organizations.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_swr_organizations" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the organization.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the organization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `organizations` - All organizations that match the filter parameters.
  The [organizations](#Organizations) structure is documented below.

<a name="Organizations"></a>
The `organizations` block supports:

* `id` - The ID of the organization.

* `name` - The name of the organization.

* `creator` - The creator of the organization.

* `permission` - The permission of organization.

* `access_user_count` - The number of users with permissions in this organization.

* `repo_count` - The number of images in this organization.
