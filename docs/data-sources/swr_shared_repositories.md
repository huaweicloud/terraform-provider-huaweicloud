---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_shared_repositories"
description: |-
  Use this data source to get the list of shared repositories.
---

# huaweicloud_swr_shared_repositories

Use this data source to get the list of shared repositories.

## Example Usage

```hcl
data "huaweicloud_swr_shared_repositories" "test" {
  center = "self"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `center` - (Required, String) Specifies the center of the repository.
  + **self:** - Indicates my shared images.
  + **thirdparty:** - Indicates images shared by others with me.

* `organization` - (Optional, String) Specifies the name of the organization (namespace) the repository belongs.

* `name` - (Optional, String) Specifies the name of the repository.

* `domain_name` - (Optional, String) Specifies the account name of the repository owner.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repositories` - All repositories that match the filter parameters.
  The [repositories](#swr_shared_repositories) structure is documented below.

<a name="swr_shared_repositories"></a>
The `repositories` block supports:

* `organization` - The name of the organization (namespace) the repository belongs.

* `name` - The name of the repository.

* `domain_name` - The account name of the repository owner.

* `status` - Whether the repository sharing has expired.

* `category` - The category of the repository.

* `is_public` - Whether the repository is public.

* `description` - The description of the repository.

* `size` - The size of the repository in byte.

* `num_download` - The number of downloads from the repository.

* `num_images` - The number of images in the repository.

* `path` - The image address for docker pull.

* `internal_path` - The intra-cluster image address for docker pull.

* `tags` - Image tag list of the repository.

* `created_at` - The creation time of the repository.

* `updated_at` - The update time of the repository.
