---
subcategory: "KooGallery"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_koogallery_assets"
description: ""
---

# huaweicloud_koogallery_assets

Use this data source to get available HuaweiCloud KooGallery assets.

## Example Usage

```hcl
variable "asset_id" {}

data "huaweicloud_koogallery_assets" "flavor" {
  asset_id      = var.asset_id
  deployed_type = "software_package"
  asset_version = "V1.0"
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region to filter the assets.

* `asset_id` - (Required, String) Specifies the asset id to filter the assets.

* `deployed_type` - (Required, String) Specifies the deployed type. Value: software_package, image.

* `asset_version` - (Optional, String) Specifies the asset version to filter the assets.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assets` - Indicates the assets information.
  The [assets](#koogallery_assets_object) structure is documented below.

<a name="koogallery_assets_object"></a>
The `assets` block contains:

* `asset_id` - The ID of the asset.
* `deployed_type` - The deployed type of the asset.
* `version` - The version of the asset.
* `version_id` - The version id of the asset.
* `region` - The region of the asset.
* `image_deployed_object` - The image deployed object information.
  The [image_deployed_object](#koogallery_assets_imsdeployed_object) structure is documented below.
* `software_pkg_deployed_object` - The software pkg deployed object information.
  The [software_pkg_deployed_object](#koogallery_assets_swdeployed_object) structure is documented below.

<a name="koogallery_assets_imsdeployed_object"></a>
The `image_deployed_object` block contains:

* `image_id` - The ID of the image asset.
* `image_name` - The name of the image asset.
* `os_type` - The os type of the image asset.
* `create_time` - The create time of the image asset.
* `architecture` - The architecture of the image asset.

<a name="koogallery_assets_swdeployed_object"></a>
The `software_pkg_deployed_object` block contains:

* `package_name` - The package name of the asset.
* `internal_path` - The internal path of the asset.
* `checksum` - The checksum of the asset.
