---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_download_secret_backup"
description: |-
  Downloads a secret backup from HuaweiCloud DEW CSMS service.
---

# huaweicloud_csms_download_secret_backup

Downloads a secret backup from HuaweiCloud DEW CSMS service.

-> This resource is a one-time action resource. Deleting this resource will not affect the actual secret backup,
  but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_csms_download_secret_backup" "test" {
  secret_name = "my-secret"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `secret_name` - (Required, String, NonUpdatable) Specifies the name of the secret to download backup for.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the `secret_name`.

* `secret_blob` - The secret backup file contains all secret version information. The backup file is encrypted and
  encoded and cannot be directly read.
