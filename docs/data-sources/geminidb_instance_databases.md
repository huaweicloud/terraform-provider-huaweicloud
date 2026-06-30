---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instance_databases"
description: |-
  Use this data source to query the databases of a GeminiDB Redis instance within HuaweiCloud.
---

# huaweicloud_geminidb_instance_databases

Use this data source to query the databases of a GeminiDB Redis instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gemini_db_instance_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instance databases.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the GeminiDB Redis instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The list of database names in the Redis instance.  
  Only databases that contain data are returned.
