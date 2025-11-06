---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_app_download_access_key"
description: |-
  Downloads a CPCS application access key resource within HuaweiCloud.
---

# huaweicloud_cpcs_app_download_access_key

Downloads a CPCS application access key resource within HuaweiCloud.

-> Each key can only be downloaded once. Downloading it repeatedly will result in an error.
 Please keep your key information safe.
 Currently, this resource is valid only in cn-north-9 region.

-> This resource is only a one-time action resource to download a CPCS application access key. Deleting this resource
will not clear the corresponding access key, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "app_id" {}
variable "access_key_id" {}

resource "huaweicloud_cpcs_app_download_access_key" "test" {
  app_id        = var.app_id
  access_key_id = var.access_key_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to download the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `app_id` - (Required, String, NonUpdatable) Specifies the application ID to which the access key belongs.

* `access_key_id` - (Required, String, NonUpdatable) Specifies the ID of the access key to download.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `access_key_id`).

* `access_key` - The access key (AK) of the downloaded key.

* `secret_key` - The secret key (SK) of the downloaded key.

* `key_name` - The name of the access key.

* `is_imported` - Whether the access key is imported.
