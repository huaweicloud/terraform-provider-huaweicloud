---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_domain_overviews"
description: |-
  Use this data source to get the list of SWR domain information.
---

# huaweicloud_swr_domain_overviews

Use this data source to get the list of SWR domain information.

## Example Usage

```hcl
data "huaweicloud_swr_domain_overviews" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namspace_num` - The number of namespaces of the tenant.

* `repo_num` - The number of repositories of the tenant.

* `image_num` - The number of images of the tenant.

* `store_space` - The storage size of the tenant.

* `downflow_size` - The download traffic of the tenant.

* `domain_id` - The domain ID.

* `domain_name` - The domain name.
