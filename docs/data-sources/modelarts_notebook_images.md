---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook_images"
description: |-
  Use this data source to get a list of ModelArts notebook images.
---

# huaweicloud_modelarts_notebook_images

Use this data source to get a list of ModelArts notebook images.

## Example Usage

### Query all notebook images without any filter

```hcl
data "huaweicloud_modelarts_notebook_images" "test" {}
```

### Query the notebook images using type filter

```hcl
data "huaweicloud_modelarts_notebook_images" "test" {
  type = "BUILD_IN"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the images are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the image to be queried.

* `type` - (Optional, String) Specifies the type of the image to be queried.  
  The valid values are as follows:
  + **BUILD_IN**: The system built-in image.
  + **DEDICATED**: User-saved images.

  Defaults to **BUILD_IN**.

* `namespace` - (Optional, String) Specifies the name of the namespace to which images belong.

* `workspace_id` - (Optional, String) Specifies the workspace ID to which images belong.

* `cpu_arch` - (Optional, String) Specifies the CPU architecture of the image to be queried.
  The valid values are as follows:
  + **x86_64**
  + **aarch64**

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `images` - The list of images that match the filter parameters.  
  The [images](#modelarts_notebook_images_attr) structure is documented below.

<a name="modelarts_notebook_images_attr"></a>
The `images` block supports:

* `id` - The ID of the image.

* `name` - The name of the image.

* `namespace` - The name of the namespace to which the image belongs.

* `workspace_id` - The workspace ID to which the image belongs.

* `resource_categories` - The supported flavors of the image.
  + **CPU**
  + **GPU**
  + **ASCEND**

* `dev_services` - The supported services of the image.  
  + **NOTEBOOK**
  + **SSH**

* `service_type` - The supported service type of the image.
  + **COMMON**
  + **INFERENCE**
  + **TRAIN**
  + **DEV**
  + **UNKNOWN**

* `show_name` - The name of the image to be displayed.

* `swr_path` - The storage path of the image.

* `type` - The type of the image.
  + **BUILD_IN**
  + **DEDICATED**

* `description` - The description of the image.

* `cpu_arch` - The CPU architecture of the image.
  + **x86_64**
  + **aarch64**

* `status` - The status of the image.
  + **ERROR**
  + **ACTIVE**
