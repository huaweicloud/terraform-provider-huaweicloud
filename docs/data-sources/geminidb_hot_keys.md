---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_hot_keys"
description: |-
  Use this data source to get the list of GeminiDB Redis instance hot keys.
---

# huaweicloud_geminidb_hot_keys

Use this data source to get the list of GeminiDB Redis instance hot keys.

-> This data source only supports GeminiDB Redis instances.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_hot_keys" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keys` - The list of GeminiDB Redis instance hot keys.
  The [keys](#keys_struct) structure is documented below.

<a name="keys_struct"></a>
The `keys` block supports:

* `name` - The name of the hot key.

* `type` - The type of the hot key.
  + **string**
  + **hash**
  + **list**
  + **zset**
  + **set**
  + **exhash**
  + **stream**

* `command` - The command of the hot key.

* `qps` - The QPS of the hot key.

* `db_id` - The DB ID to which the hot key belongs.
