---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_all_repositories"
description: |-
  Use this data source to get the list of SWR enterprise all repositories.
---

# huaweicloud_swr_enterprise_all_repositories

Use this data source to get the list of SWR enterprise all repositories in the current project,
excluding repositories in shared Enterprise Edition instances.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_swr_enterprise_all_repositories" "test" {}
```

### Filter By Repository Name

```hcl
variable "repo_name" {}

data "huaweicloud_swr_enterprise_all_repositories" "test" {
  name = var.repo_name
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the repository name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repositories` - Indicates the repository list.

  The [repositories](#repositories_struct) structure is documented below.

<a name="repositories_struct"></a>
The `repositories` block supports:

* `id` - Indicates the repository ID.

* `name` - Indicates the repository name.

* `namespace_name` - Indicates the namespace name.

* `namespace_id` - Indicates the namespace ID.

* `tag_count` - Indicates the number of artifact tags in a repository.

* `pull_count` - Indicates the total number of downloads.

* `artifact_count` - Indicates the total number of artifact packages.

* `description` - Indicates the description.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

* `instance_id` - Indicates the ID of an SWR Enterprise Edition instance.

* `instance_name` - Indicates the name of an SWR Enterprise Edition instance.

* `resource_urn` - Indicates the resource URN.
