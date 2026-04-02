---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_database_character_set"
description: |-
  Use this data source to get the list of TaurusDB database character sets.
---

# huaweicloud_taurusdb_database_character_set

Use this data source to get the list of TaurusDB database character sets.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_database_character_set" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `charsets` - Indicates the list of database character sets.
