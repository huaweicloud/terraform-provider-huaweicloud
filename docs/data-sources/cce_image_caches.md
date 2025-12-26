---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_image_caches"
description: |-
  Use this data source to get the image caches.
---

# huaweicloud_cce_image_caches

Use this data source to get the image caches.

## Example Usage

```hcl
data "huaweicloud_cce_image_caches" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `image_caches` - The image caches data in cce cluster.

  The [image_caches](#image_caches_struct) structure is documented below.

<a name="image_caches_struct"></a>
The `image_caches` block supports:

* `name` - The image caches name.

* `id` - The image caches id.

* `created_at` - The image caches create time.

* `images` - The image list in image caches.

* `image_cache_size` - The evs size in image caches.

* `retention_days` - The retention days in image caches.

* `status` - The image caches status.

* `message` - The messages in image caches.
