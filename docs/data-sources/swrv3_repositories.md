---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swrv3_repositories"
description: |-
  Use this data source to get the list of SWR repositories.
---

# huaweicloud_swrv3_repositories

Use this data source to get the list of SWR repositories.

## Example Usage

```hcl
data "huaweicloud_swrv3_repositories" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `organization` - (Optional, String) Specifies the organization name.
  Enter 1 to 64 characters, starting with a lowercase letter and ending with a lowercase letter or digit.
  Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed. Periods, underscores,
  and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.

* `name` - (Optional, String) Specifies the image repository name. Enter 1 to 128 characters. It must start and end with
  a lowercase letter or digit. Only lowercase letters, digits, periods (.), slashes (/), underscores (_), and hyphens (-)
  are allowed. Periods, slashes, underscores, and hyphens cannot be placed next to each other. A maximum of two
  consecutive underscores are allowed. Replace a slash (/) with a dollar sign ($) before you send the request.

* `category` - (Optional, String) Specifies the repository type.
  The value can be **app_server**, **linux**, **framework_app**, **database**, **lang**, **windows**, **arm**, or other.

* `is_public` - (Optional, String) Specifies whether a repository is public. Valid value is **true** and **false**.
  Default to return all.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repos` - Indicates the repositories.

  The [repos](#repos_struct) structure is documented below.

<a name="repos_struct"></a>
The `repos` block supports:

* `id` - Indicates the repository ID.

* `name` - Indicates the repository name.

* `size` - Indicates the repository size.

* `namespace_name` - Indicates the name of the organization that a repository belongs to.

* `num_download` - Indicates the repository downloads.

* `status` - Indicates whether the image shared by others has expired.
  This field is not returned when the repository information of the current user is obtained.

* `category` - Indicates the repository type.

* `is_public` - Indicates whether a repository is a public repository.

* `num_images` - Indicates the number of images in a repository.

* `description` - Indicates the repository description.

* `created_at` - Indicates the time when a repository was created. It is the UTC standard time.
  Users need to calculate the offset based on the local time, for example, UTC+8:00 for the East 8th Time Zone.

* `updated_at` - Indicates the time when a repository was updated. It is the UTC standard time.
  Users need to calculate the offset based on the local time, for example, UTC+8:00 for the East 8th Time Zone.
