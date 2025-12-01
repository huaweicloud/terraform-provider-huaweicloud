---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_upload_command"
description: |-
  Use this data source to get the image upload command of HSS within HuaweiCloud.
---

# huaweicloud_hss_image_upload_command

Use this data source to get the image upload command of HSS within HuaweiCloud.

## Example Usage

```hcl
variable "registry_addr" {}
variable "namespace" {}
variable "username" {}
variable "password" {}

data "huaweicloud_hss_image_upload_command" "test" {
  registry_addr = var.registry_addr
  namespace     = var.namespace
  username      = var.username
  password      = var.password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `registry_addr` - (Required, String) Specifies the image repository address.

* `namespace` - (Required, String) Specifies the image repository project, which is used to specify the image repository
  directory that the scan component is to be uploaded to.

* `username` - (Required, String) Specifies the username for logging in to the image repository.

* `password` - (Required, String) Specifies the password for logging in to the image repository. The password is used
  only to generate commands. HSS does not store your image repository password.
  This field is sensitive.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `image_command` - The command for uploading an image.

* `images_download_url` - The image download link.

* `swr_image_pull_command` - The command for obtaining an image from SWR.
