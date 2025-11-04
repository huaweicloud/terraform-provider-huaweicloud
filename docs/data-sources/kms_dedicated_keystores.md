---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_dedicated_keystores"
description: |-
  Use this datasource to get the list of dedicated keystores.
---

# huaweicloud_kms_dedicated_keystores

Use this datasource to get the list of dedicated keystores.

## Example Usage

```hcl
data "huaweicloud_kms_dedicated_keystores" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keystores` - The list of the keystores.

  The [keystores](#keystores_struct) structure is documented below.

<a name="keystores_struct"></a>
The `keystores` block supports:

* `keystore_id` - The keystore ID.

* `domain_id` - The user domain ID.

* `keystore_alias` - The keystore alias.

* `keystore_type` - The keystore type.

* `hsm_cluster_id` - The DHSM cluster ID.

* `cluster_id` - The cluster ID.

* `create_time` - The creation time of the keystore.
