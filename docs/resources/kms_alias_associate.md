---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_alias_associate"
description: |-
  Manages a KMS alias associate resource within HuaweiCloud
---

# huaweicloud_kms_alias_associate

Manages a KMS alias associate resource within HuaweiCloud.

  -> This resource relies on the alias created by the huaweicloud_kms_alias resource to associate with the target key.

  -> When the alias is associated with the target key, the alias on the original key will be destroyed.

## Example Usage

```hcl
variable "target_key_id" {}
variable "alias" {}

resource "huaweicloud_kms_alias_associate" "test" {
  target_key_id = var.target_key_id
  alias         = var.alias
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `target_key_id` - (Required, String, NonUpdatable) Specifies the target key ID used to bind the alias.

* `alias` - (Required, String, NonUpdatable) Specifies the alias of the key, it can only be prefixed with **alias/**.

  -> And the alias must be an alias already used by other keys.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is **target_key_id?alias**.

* `alias_urn` - The alias resource locator.

* `create_time` - The creation time of the alias.

* `update_time` - The update time of the alias.

## Import

The KMS alias can be imported using `target_key_id` and `alias`, separated by a question mark (?), e.g.

```bash
$ terraform import huaweicloud_kms_alias_associate.test <target_key_id>?<alias>
```
