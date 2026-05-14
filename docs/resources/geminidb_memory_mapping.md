---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_memory_mapping"
description: |-
  Manages a memory mapping resource within HuaweiCloud.
---

# huaweicloud_geminidb_memory_mapping

Manages a memory mapping resource within HuaweiCloud.

-> This resource only supports primary/standby GeminiDB Redis instance.

## Example Usage

```hcl
var "source_instance_id" {}
var "target_instance_id" {}

resource "huaweicloud_geminidb_memory_mapping" "test" {
  source_instance_id = var.source_instance_id
  target_instance_id = var.target_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `source_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the source instance.
  Currently, RDS for MySQL and TaurusDB instances are supported.

* `target_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the destination instance.
  Currently, only primary/standby GeminiDB Redis instance is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - The memory mapping name.

* `source_instance_name` - The source instance name.

* `target_instance_name` - The target instance name.

* `status` - The memory mapping status.

* `created` - The memory mapping creation time.

* `updated` - The memory mapping update time.

* `rule_count` - The number of rules associated with the memory mapping.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_geminidb_memory_mapping.test <id>
```
