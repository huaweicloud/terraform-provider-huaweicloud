---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swrv3_shared_repositories"
description: |-
  Use this data source to get the list of SWR shared repositories.
---

# huaweicloud_swrv3_shared_repositories

Use this data source to get the list of SWR shared repositories.

## Example Usage

```hcl
data "huaweicloud_swrv3_shared_repositories" "test" {
  shared_by = "self"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `shared_by` - (Required, String) Specifies the sharing mode.
  The value can only be as follows:
  + **self**: images shared by you to others
  + **thirdparty**: images shared with you by others

* `name` - (Optional, String) Specifies the image repository name. Enter 1 to 128 characters. It must start and end with
  a lowercase letter or digit. Only lowercase letters, digits, periods (.), slashes (/), underscores (_), and hyphens (-)
  are allowed. Periods, slashes, underscores, and hyphens cannot be placed next to each other. A maximum of two
  consecutive underscores are allowed. Replace a slash (/) with a dollar sign ($) before you send the request.

* `organization` - (Optional, String) Specifies the organization name. Enter 1 to 64 characters, starting with a
  lowercase letter and ending with a lowercase letter or digit. Only lowercase letters, digits, periods (.),
  underscores (_), and hyphens (-) are allowed. Periods, underscores, and hyphens cannot be placed next to each other.
  A maximum of two consecutive underscores are allowed.

* `status` - (Optional, String) Specifies whether the sharing has expired.
  The value can only be as follows:
  + **false**: The sharing has expired.
  + **true**: The sharing is valid.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repos` - Indicates the repositories.
  The [repos](#attrblock--repos) structure is documented below.

<a name="attrblock--repos"></a>
The `repos` block supports:

* `id` - Indicates the repository ID.

* `name` - Indicates the repository name.

* `size` - Indicates the repository size.

* `category` - Indicates the repository type.

* `description` - Indicates the repository description.

* `is_public` - Indicates whether a repository is a public repository.

* `num_download` - Indicates the repository downloads.

* `num_images` - Indicates the number of images in a repository.

* `organization` - Indicates the name of the organization that a repository belongs to.

* `status` - Indicates whether the image shared by others has expired.

* `created_at` - Indicates the time when a repository was created. It is the UTC standard time.

* `updated_at` - Indicates the time when a repository was updated. It is the UTC standard time.
