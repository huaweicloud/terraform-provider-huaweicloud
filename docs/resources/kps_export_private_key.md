---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_export_private_key"
description: |-
  Manages a KPS export private key resource within HuaweiCloud.
---

# huaweicloud_kps_export_private_key

Manages a KPS export private key resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the exported key,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_kps_keypair" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the SSH keypair name to export the private key from.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `private_key` - The private key of the keypair. This is marked as sensitive and will not be
  displayed in logs.
