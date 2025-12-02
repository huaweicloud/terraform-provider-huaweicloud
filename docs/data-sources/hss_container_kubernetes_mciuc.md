---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_mciuc"
description: |-
  Use this data source to get the container kubernetes multi cloud image upload command of HSS within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_mciuc

Use this data source to get the container kubernetes multi cloud image upload command of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_mciuc" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `image_repo` - (Required, String) Specifies the image repository address.

* `organization` - (Required, String) Specifies the organization name.

* `username` - (Required, String) Specifies the username.

* `password` - (Required, String) Specifies the password.

* `plug_type` - (Optional, String) Specifies the plug-in type. Valid values are:
  + **docker**: Docker plug-in image
  + **agent**: Hostguard image

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `image_command` - The command for uploading an image.

* `secret_command` - The installation credential command.

* `images_download_url` - The image download link.
