---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_custom_images"
description: |-
  Use this data source to get custom images of CPH phone.
---

# huaweicloud_cph_phone_custom_images

Use this data source to get custom images of CPH phone.

## Example Usage

```hcl
data "huaweicloud_cph_phone_custom_images" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the image type. The valid value can be **public**, **private** or **share**.

* `status` - (Optional, String) Specifies the image status.
  The valid value can be **0** (creating), **-1** (production) or **-2** (create failed).

* `image_id` - (Optional, String) Specifies the image ID.

* `name` - (Optional, String) Specifies the image name.

* `create_since` - (Optional, String) Specifies the image creation since time.

* `create_until` - (Optional, String) Specifies the image creation until time.

* `src_project_id` - (Optional, String) Specifies the project ID of the share image account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - The image list.

  The [images](#images_struct) structure is documented below.

<a name="images_struct"></a>
The `images` block supports:

* `name` - The image name.

* `size` - The image size, the unit is byte.

* `project_id` - The project ID of the image.

* `id` - The image ID.

* `version` - The image AOSP version.

* `status` - The image status.

* `src_project_id` - The project ID of the share image account.

* `domain_id` - The domain ID to which the image belongs.

* `updated_at` - The image update time.

* `created_at` - The image creation time.

* `type` - The image type.
