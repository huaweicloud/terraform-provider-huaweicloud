---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_alias"
description: |-
  Manages a KMS alias resource within HuaweiCloud
---

# huaweicloud_kms_alias

Manages a KMS alias resource within HuaweiCloud.

## Example Usage

```hcl
variable "key_id" {}
variable "alias" {}

resource "huaweicloud_kms_alias" "test" {
  key_id = var.key_id
  alias  = var.alias
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID used to bind the alias.

* `alias` - (Required, String, NonUpdatable) Specifies the alias of the key, it can only be prefixed with **alias/**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is **key_id?alias**.

* `alias_urn` - The alias resource locator.

* `create_time` - The creation time of the alias.

* `update_time` - The update time of the alias.

## Import

The KMS alias can be imported using `key_id` and `alias`, separated by a question mark (?), e.g.

```bash
$ terraform import huaweicloud_kms_alias.test <key_id>?<alias>
```
