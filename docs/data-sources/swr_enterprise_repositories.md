---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_repositories"
description: |-
  Use this data source to get the list of SWR enterprise repositories.
---

# huaweicloud_swr_enterprise_repositories

Use this data source to get the list of SWR enterprise repositories.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_id` - (Optional, String) Specifies the namespace ID.

* `order_column` - (Optional, String) Specifies the order column.
  Values can be **created_at** or **updated_at**. Default to **created_at**.

* `order_type` - (Optional, String) Specifies the order type. Values can be **desc** or **asc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `repositories` - Indicates the repositories.

  The [repositories](#repositories_struct) structure is documented below.

* `total` - Indicates the total count.

<a name="repositories_struct"></a>
The `repositories` block supports:

* `id` - Indicates the repository ID.

* `name` - Indicates the repository name.

* `namespace_id` - Indicates the namespace ID.

* `namespace_name` - Indicates the namespace ID.

* `pull_count` - Indicates the count of pull.

* `artifact_count` - Indicates the count of artifact.

* `description` - Indicates the description.

* `tag_count` - Indicates the count of tags.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
