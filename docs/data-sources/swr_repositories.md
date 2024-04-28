---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_repositories"
description: ""
---

# huaweicloud_swr_repositories

Use this data source to get the list of repositories.

## Example Usage

```hcl
data "huaweicloud_swr_repositories" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `organization` - (Optional, String) Specifies the name of the organization (namespace) the repository belongs.

* `name` - (Optional, String) Specifies the name of the repository.

* `category` - (Optional, String) Specifies the category of the repository. The value can be **app_server**,
  **linux**, **framework_app**, **database**, **lang**, **other**, **windows**, **arm**.

* `is_public` - (Optional, String) Specifies whether the repository is public. Default is false.
  + **true** - Indicates the repository is public.
  + **false** - Indicates the repository is private.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repositories` - All repositories that match the filter parameters.
  The [repositories](#attrblock--repositories) structure is documented below.

<a name="attrblock--repositories"></a>
The `repositories` block supports:

* `organization` - The name of the organization (namespace) the repository belongs.

* `name` - The name of the repository.

* `category` - The category of the repository.

* `is_public` - Whether the repository is public.

* `description` - The description of the repository.

* `size` - The size of the repository in byte.

* `num_download` - The number of downloads from the repository.

* `num_images` - The number of images in the repository.

* `path` - The image address for docker pull.

* `internal_path` - The intra-cluster image address for docker pull.

* `tags` - Image tag list of the repository.

* `status` - Whether this repository is shared with others.

* `total_range` - The total number of the repository.

* `created_at` - The creation time of the repository.

* `updated_at` - The update time of the repository.
