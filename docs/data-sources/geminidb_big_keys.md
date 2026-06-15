---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_big_keys"
description: |-
  Use this data source to get the list of GeminiDB Redis instance big keys.
---

# huaweicloud_geminidb_big_keys

Use this data source to get the list of GeminiDB Redis instance big keys.

-> This data source only supports GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_big_keys" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `key_types` - (Optional, List) Specifies the type of the big key.
  The valid values are as follows:
  + **string**
  + **hash**
  + **zset**
  + **set**
  + **list**
  + **exhash**
  + **stream**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keys` - The list of GeminiDB Redis instance big keys.
  The [keys](#keys_struct) structure is documented below.

<a name="keys_struct"></a>
The `keys` block supports:

* `db_id` - The DB ID to which the big key belongs.

* `key_type` - The type of the big key.

* `key_name` - The name of the big key.

* `key_size` - The length of the big key.
