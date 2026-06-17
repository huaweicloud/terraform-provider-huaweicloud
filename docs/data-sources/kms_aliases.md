---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_aliases"
description: |-
  Use this data source to query the list of KMS aliases within HuaweiCloud.
---

# huaweicloud_kms_aliases

Use this data source to query the list of KMS aliases within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_kms_aliases" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `key_id` - (Optional, String) Specifies the key ID used to query the alias.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `aliases` - The list of key aliases.
  The [aliases](#aliases_struct) structure is documented below.

<a name="aliases_struct"></a>
The `aliases` block supports:

* `domain_id` - The ID of the account to which the alias belongs.

* `key_id` - The key ID.

* `alias` - The alias of the key.

* `alias_urn` - The alias resource locator.

* `create_time` - The creation time of the alias.

* `update_time` - The update time of the alias.
